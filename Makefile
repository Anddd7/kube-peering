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
example_tcp_proxy:
	cd example/tcp/proxy && go run . $(args) && cd -
example_tcp_vpn:
	cd example/tcp/vpn && go run . $(args) && cd -

example_http_app:
	cd example/http/app && go run . $(args) && cd -
example_http_client:
	cd example/http/client && go run . $(args) && cd -
example_http_proxy:
	cd example/http/proxy && go run . $(args) && cd -
example_http_vpn:
	cd example/http/vpn && go run . $(args) && cd -

example_https_app:
	cd example/https/app && go run . $(args) && cd -
example_https_client:
	cd example/https/client && go run . $(args) && cd -
example_https_proxy:
	cd example/https/proxy && go run . $(args) && cd -
example_https_vpn:
	cd example/https/vpn && go run . $(args) && cd -