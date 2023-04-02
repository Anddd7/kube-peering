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

# -------------------------------- #

example_tcp_app:
	cd example/tcp/app && go run . $(args) && cd -
example_tcp_client:
	cd example/tcp/client && go run . $(args) && cd -
example_tcp_transit:
	cd example/tcp/transit && go run . $(args) && cd -

example_http_app:
	cd example/http/app && go run . $(args) && cd -
example_http_client:
	cd example/http/client && go run . $(args) && cd -
example_http_transit:
	cd example/http/transit && go run . $(args) && cd -

example_https_app:
	cd example/https/app && go run . $(args) && cd -
example_https_client:
	cd example/https/client && go run . $(args) && cd -
example_https_transit:
	cd example/https/transit && go run . $(args) && cd -