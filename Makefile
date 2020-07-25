include .env

MAKEFLAGS += --silent
PRJ_PATH=$(PWD)
PRJ_NAME=$(shell basename "$(PWD)")
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test -v
GOMOD=$(GOCMD) mod
TESTPATH=$(PRJ_PATH)/internal/pkg/test
IMAGE_TAG=latest

configs.app.yaml: configs/app-dev.yaml
	cp configs/app-dev.yaml configs/app.yaml 

run.ws:
	@export ENV=$(ENV)
	$(GORUN) main.go ws

run.client:
	$(GORUN) main.go client


build: main.go
	$(GOBUILD) -o $(PRJ_NAME) main.go

mod.vendor:
	$(GOMOD) tidy && \
	$(GOMOD) vendor

go.proto:
	protoc -I $(PRJ_PATH)/pkg/echoproto -I $(PRJ_PATH)/vendor  --go_out=$(PRJ_PATH)/pkg/echoproto \
	--go_opt=paths=source_relative $(PRJ_PATH)/pkg/echoproto/*.proto

go.lint:
	golint $(PRJ_PATH)/internal/pkg/app/...

test:
	$(GOTEST) $(TESTPATH)/$(test)_test -count=1

test.all:
	$(GOTEST) $(TESTPATH)/... -count=1

test.pkg:
	$(GOTEST) $(PRJ_PATH)/pkg/...

bench.proto: $(PRJ_PATH)/pkg/echoproto/proto_test.go
	$(GOTEST) $(PRJ_PATH)/pkg/echoproto -run=None -bench=. --benchmem

build.image:
	docker build -t=$(DOCKEHUB)/whaleecho:$(IMAGE_TAG) -f $(PRJ_PATH)/deployments/docker/Dockerfile .

push.image:
	docker push $(DOCKEHUB)/whaleecho

start.containers: $(PRJ_PATH)/deployments/docker/docker-compose.yaml
	docker-compose -p $(PRJ_NAME) -f $(PRJ_PATH)/deployments/docker/docker-compose.yaml up

teardown.containers: $(PRJ_PATH)/deployments/docker/docker-compose.yaml
	docker-compose -p $(PRJ_NAME) -f $(PRJ_PATH)/deployments/docker/docker-compose.yaml down
	docker container prune --force

push.repo:
	git push origin HEAD

setup.env:
	. $(PRJ_PATH)/scripts/setup.sh

ci: go.lint build test.all push.repo build.image push.image

build.nats:
	docker run -d --name nats -p 4222:4222 -p 6222:6222 -p 8222:8222 nats 
# cd: 

# cicd: ci cd