include .env

MAKEFLAGS += --silent
PRJ_NAME=$(shell basename "$(PWD)")
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
PID := ./tmp/.$(PROJECTNAME).pid
TESTPATH=$(PWD)/internal/pkg/test

configs.app.yaml: configs/app-dev.yaml
	cp configs/app-dev.yaml configs/app.yaml 

run.ws:
	@export ENV=$(ENV)
	$(GORUN) main.go ws

build: main.go
	$(GOBUILD) -o $(PRJ_NAME) main.go

mod.vendor:
	$(GOMOD) tidy && \
	$(GOMOD) vendor

test:
	$(GOTEST) -v $(TESTPATH)/$(test)_test -count=1
