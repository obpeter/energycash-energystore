# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=energystore
DOCKER=docker

GOPATH := ${PWD}/..:${GOPATH}
export GOPATH

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags="-s -w"
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

docker-clean:
	$(DOCKER) rmi ghcr.io/vfeeg-development/energy-store:latest

docker: docker-clean
	$(DOCKER) build -t ghcr.io/vfeeg-development/energy-store:latest .

push: docker
	$(DOCKER) push ghcr.io/vfeeg-development/energy-store:latest