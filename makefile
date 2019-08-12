build:
	mkdir -p ./build && go build -ldflags "-w -s" -o ./octopus ./main.go ./ws.go ./http.go ./router.go
build-linux: 
	CGO_ENABLED=no GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ./octopus ./main.go ./ws.go ./http.go ./router.go
build-mac:
	CGO_ENABLED=no go build -ldflags "-w -s" -o ./octopus-mac ./main.go ./ws.go ./http.go ./router.go