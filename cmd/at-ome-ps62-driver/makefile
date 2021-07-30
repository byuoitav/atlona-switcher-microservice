NAME := atlona-switcher-microservice
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}
DOCKER_URL := docker.pkg.github.com
DOCKER_PKG := ${DOCKER_URL}/${OWNER}/${NAME}
BRANCH:= $(shell git rev-parse --abbrev-ref HEAD)

# go
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
VENDOR=gvt fetch -branch $(BRANCH)

# docker
DOCKER=docker
DOCKER_BUILD=$(DOCKER) build
DOCKER_LOGIN=$(DOCKER) login -u $(UNAME) -p $(PASS)
DOCKER_PUSH=$(DOCKER) push
DOCKER_FILE=dockerfile

UNAME=$(shell echo $(DOCKER_USERNAME))
EMAIL=$(shell echo $(DOCKER_EMAIL))
PASS=$(shell echo $(DOCKER_PASSWORD))

build: deps
	@mkdir -p dist
	@env GOOS=linux CGO_ENABLED=0 $(GOBUILD) -o ./dist/$(NAME)-linux-amd64 -v
	@env GOOS=linux GOARCH=arm $(GOBUILD) -o ./dist/$(NAME)-linux-arm -v
	
test: 
	@$(GOTEST) -v -race $(go list ./... | grep -v /vendor/) 

clean: 
	@$(GOCLEAN)
	@rm -rf dist/

deps: 
	@$(GOGET) -d -v

docker: clean build
	@echo Building amd64 container
	@$(DOCKER_BUILD) --build-arg NAME=$(NAME)-linux-amd64 -f $(DOCKER_FILE) -t $(DOCKER_PKG)/$(NAME):$(BRANCH) dist

	@echo Building arm container
	@$(DOCKER_BUILD) --build-arg NAME=$(NAME)-linux-arm -f $(DOCKER_FILE) -t $(DOCKER_PKG)/rpi-$(NAME):$(BRANCH) dist

deploy: docker
	@echo logging in to dockerhub...
	@$(DOCKER_LOGIN)

	@$(DOCKER_PUSH) $(ORG)/rpi-$(NAME):$(BRANCH)
	@$(DOCKER_PUSH) $(ORG)/$(NAME):$(BRANCH)
