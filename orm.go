package main

import (
	"context"
	"net/http"
	"os"

	"github.com/graphql-services/go-saga/graphqlorm"
	"github.com/machinebox/graphql"
	"golang.org/x/oauth2/clientcredentials"
)

// var client *graphqlorm.ORMClient

var clientcredentialsClient *http.Client

func getClientcredentialsClient() *http.Client {
	if clientcredentialsClient == nil {
		clientID := os.Getenv("OIDC_CLIENT_ID")
		if clientID == "" {
			panic("Missing 'OIDC_CLIENT_ID' environment variable")
		}
		clientSecret := os.Getenv("OIDC_CLIENT_SECRET")
		if clientSecret == "" {
			panic("Missing 'OIDC_CLIENT_SECRET' environment variable")
		}
		tokenURL := os.Getenv("OIDC_TOKEN_URL")
		if tokenURL == "" {
			panic("Missing OIDC_TOKEN_URL environment variable")
		}
		conf := &clientcredentials.Config{
			TokenURL:     tokenURL,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}
		clientcredentialsClient = conf.Client(context.Background())
	}
	return clientcredentialsClient
}

func GetORMClient() *graphqlorm.ORMClient {
	client := getClientcredentialsClient()
	URL := os.Getenv("ORM_URL")
	if URL == "" {
		panic("Missing ORM_URL environment variable")
	}
	return graphqlorm.NewClient(URL, graphql.WithHTTPClient(client))

}
