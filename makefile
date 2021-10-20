build:
	mkdir -p ./build && go build -ldflags "-w -s" -o ./build/octopus ./main.go ./ws.go ./http.go ./router.go
build-linux: 
	CGO_ENABLED=no GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ./build/octopus ./main.go ./ws.go ./http.go ./router.go
build-mac:
	CGO_ENABLED=no go build -ldflags "-w -s" -o ./build/octopus-mac ./main.go ./ws.go ./http.go ./router.go

.PHONY: cover build build-linux build-mac
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html