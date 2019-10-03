vendors:
	go mod vendor
	go mod download

test:	
	go test -v -mod=vendor ./cmd/*.go

build:
	go build -mod=vendor -ldflags="-w -s" -o .bin/multi-tag-searcher cmd/*.go	

run:build
	.bin/multi-tag-searcher

docker-build: 
	docker build -t multi-tag-searcher .

docker-run:
	docker run -it --rm multi-tag-searcher