test:
	go test -coverprofile=c.out -covermode=atomic ./yamlutils

view:
	go tool cover -html=c.out

build:
	go build -v -o bin/yaml-parse cmd/yaml-parse/main.go

release:
	go build -v -o bin/yaml-parse \
		-ldflags="-X main.BuildMetadata=`date +'%Y%m%d%H%M%S'`.`git rev-parse --short HEAD`" \
		cmd/yaml-parse/main.go
