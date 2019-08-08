clean:
	go clean && rm -rf build/*
build: clean
	mkdir -p ./build && go build -ldflags "-w -s" -o ./build/antman ./main.go
build-linux: clean 
	CGO_ENABLED=no GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ./octopus-linux ./main.go ./ws.go ./http.go ./router.go
build-mac: clean 
	CGO_ENABLED=no go build -ldflags "-w -s" -o ./octopus-mac ./main.go ./ws.go ./http.go ./router.go
build-opcap: clean
	CGO_ENABLED=no GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ./opcap-linux ./opcap.go
build-opcap-cpp: clean
	gcc -o opcap ./opcap.cpp -lpcap
.IGNORE: clean