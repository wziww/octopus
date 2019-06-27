clean:
	go clean && rm -rf build/*
build: clean
	mkdir -p ./build && go build -ldflags "-w -s" -o ./build/antman ./main.go
build-linux: clean 
	CGO_ENABLED=no GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ./build/octopus ./main.go ./ws.go
.IGNORE: clean