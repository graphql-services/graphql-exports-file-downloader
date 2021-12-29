OWNER=graphql
IMAGE_NAME=graphql-exports-file-downloader
QNAME=$(OWNER)/$(IMAGE_NAME)

GIT_TAG=$(QNAME):$(GITHUB_SHA)
BUILD_TAG=$(QNAME):$(GITHUB_RUN_ID).$(GITHUB_SHA)
TAG=$(QNAME):`echo $(GITHUB_REF) | sed 's/refs\/heads\///' | sed 's/master/latest/;s/develop/unstable/' | sed 's/refs\/tags\///'`

lint:
	docker run -it --rm -v "$(PWD)/Dockerfile:/Dockerfile:ro" redcoolbeans/dockerlint

build:
	# go get ./...
	# gox -osarch="linux/amd64" -output="bin/devops-alpine"
	# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/binary .
	docker build -t $(GIT_TAG) .

tag:
	docker tag $(GIT_TAG) $(BUILD_TAG)
	docker tag $(GIT_TAG) $(TAG)
	
login:
	@docker login -u "$(DOCKER_USER)" -p "$(DOCKER_PASSWORD)"
push: login
	# docker push $(GIT_TAG)
	# docker push $(BUILD_TAG)
	docker push $(TAG)

build-lambda-function:
	GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build -o main *.go && zip lambda.zip main && rm main

build-local:
	go get ./...
	go build -o $(IMAGE_NAME)

deploy-local:
	make build-local
	mv $(IMAGE_NAME) /usr/local/bin/

test:
	PORT=8080 go run *.go start
