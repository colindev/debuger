all: bin/api bin/ws bin/nc bin/json-indent

bin/api: api/*.go
	go build -o bin/api api/*.go

bin/ws: ws/*.go
	go build -o bin/ws ws/*.go

bin/nc: nc/*.go
	go build -o bin/nc nc/*.go

bin/json-indent: json/*.go
	go build -o bin/json-indent json/*.go
