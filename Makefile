build-api: api/*.go
	go build -o api/api api/*.go

build-ws: ws/*.go
	go build -o ws/ws ws/*.go

all: build-api build-ws
