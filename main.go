package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/graphql-services/go-saga/graphqlorm"
	handler "github.com/jakubknejzlik/cloudevents-lambda-handler"
	"github.com/novacloudcz/graphql-orm/events"
)

func main() {
	ormClient := GetORMClient()
	h := handler.NewCloudEventsLambdaHandler(receiver(ormClient))
	h.Start()
}

type FileMetadataFile struct {
	Name string
	URL  string
}

type FileMetadata struct {
	Filename string
	Files    []FileMetadataFile
}

func receiver(ormClient *graphqlorm.ORMClient) func(e cloudevents.Event) (err error) {
	return func(e cloudevents.Event) (err error) {
		var ormEvent events.Event
		err = e.DataAs(&ormEvent)
		if err != nil {
			return
		}

		ctx := context.Background()
		if ormEvent.Entity == "Export" && ormEvent.Type == events.EventTypeCreated {
			err = handleExport(ctx, ormClient, &ormEvent)
			if err != nil {
				updateExportError(ctx, ormClient, ormEvent.EntityID, err.Error())
			}
		}

		return
	}
}

func handleExport(ctx context.Context, ormClient *graphqlorm.ORMClient, ormEvent *events.Event) (err error) {
	fmt.Println("new export", ormEvent)

	err = updateExport(ctx, ormClient, ormEvent.EntityID, "PROCESSING", 0)
	if err != nil {
		return
	}

	var metadataString string
	err = ormEvent.Change("metadata").NewValueAs(&metadataString)
	if err != nil {
		return
	}

	fmt.Println("metadata", metadataString)
	var metadata FileMetadata
	err = json.Unmarshal([]byte(metadataString), &metadata)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	buf, err = ZipFiles(metadata.Files, func(progress float32) error {
		return updateExport(ctx, ormClient, ormEvent.EntityID, "PROCESSING", progress)
	})
	if err != nil {
		return
	}
	fmt.Println("zipped", len(buf.Bytes()))

	var fileId string
	fileId, err = Upload(metadata.Filename, buf)
	if err != nil {
		return
	}

	_, err = ormClient.UpdateEntity(ctx, graphqlorm.UpdateEntityOptions{
		Entity:   ormEvent.Entity,
		EntityID: ormEvent.EntityID,
		Input: map[string]string{
			"fileId": fileId,
			"state":  "COMPLETED",
		},
	})
	return
}

func updateExport(ctx context.Context, ormClient *graphqlorm.ORMClient, exportId string, state string, progress float32) (err error) {
	_, err = ormClient.UpdateEntity(ctx, graphqlorm.UpdateEntityOptions{
		Entity:   "Export",
		EntityID: exportId,
		Input: map[string]interface{}{
			"state":    state,
			"progress": progress,
		},
	})
	return
}

func updateExportError(ctx context.Context, ormClient *graphqlorm.ORMClient, exportId string, errorDescription string) (err error) {
	_, err = ormClient.UpdateEntity(ctx, graphqlorm.UpdateEntityOptions{
		Entity:   "Export",
		EntityID: exportId,
		Input: map[string]interface{}{
			"state":            "ERROR",
			"errorDescription": errorDescription,
		},
	})
	return
}
