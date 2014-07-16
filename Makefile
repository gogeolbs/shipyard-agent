all:
	@cd ./agent && go build -o ../shipyard-agent && chmod +x ../shipyard-agent

get:
	@go get -d -v ./...
fmt:
	@go fmt ./...
test:
	@go test ./...
clean:
	@rm -rf shipyard-agent
