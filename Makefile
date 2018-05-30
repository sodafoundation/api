.PHONY: all build protoc osdsdock osdslet osdsctl docker clean

all:build

build:osdsdock osdslet osdsctl

protoc:
	cd pkg/dock/proto && protoc --go_out=plugins=grpc:. dock.proto
osdsdock:
	mkdir -p  ./build/out/bin/
	go build -o ./build/out/bin/osdsdock github.com/opensds/opensds/cmd/osdsdock

osdslet:
	mkdir -p  ./build/out/bin/
	go build -o ./build/out/bin/osdslet github.com/opensds/opensds/cmd/osdslet

osdsctl:
	mkdir -p  ./build/out/bin/
	go build -o ./build/out/bin/osdsctl github.com/opensds/opensds/osdsctl

docker:build
	cp ./build/out/bin/osdsdock ./cmd/osdsdock
	cp ./build/out/bin/osdslet ./cmd/osdslet
	docker build cmd/osdsdock -t opensdsio/opensds-dock:latest
	docker build cmd/osdslet -t opensdsio/opensds-controller:latest
test:build
	script/CI/test

clean:
	rm -rf ./build ./cmd/osdslet/osdslet ./cmd/osdsdock/osdsdock
