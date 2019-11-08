VERSION?=$(shell git describe --tags --abbrev=0)
SHORT_HASH_COMMIT=$(shell git rev-parse --short HEAD)
SERVICE?=multi-tag-searcher

DOTENV := .env
DOTENV_EXISTS := $(shell [ -f $(DOTENV) ] && echo 0 || echo 1 )

ifeq ($(DOTENV_EXISTS), 0)
	include $(DOTENV)
	export $(shell sed 's/=.*//' .env)
endif

vendors:
	go mod vendor
	go mod download

test:	
	go test -v -mod=vendor ./cmd/*.go

build:
	go build -mod=vendor -ldflags "-w -s -X main.version=${VERSION} -X main.commitHash=${SHORT_HASH_COMMIT}" -o .bin/${SERVICE} cmd/*.go

run:build
	.bin/${SERVICE}

docker-build: 
	docker build -t multi-tag-searcher .

docker-run:
	docker run -it --rm multi-tag-searcher