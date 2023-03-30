git_hooks:
	cp scripts/hooks/commit-msg .git/hooks/

clean:
	rm -rf bin

fmt:
	go fmt ./...
	gofumpt -l -w .
	go vet ./...

lint:
	golangci-lint run -v

cover:
	go test ./... -coverprofile ./bin/coverage.out

coverweb: cover
	go tool cover -html=./bin/coverage.out

check: fmt lint cover

dep:
	go mod download
