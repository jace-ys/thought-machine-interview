.PHONY: cpx-server deps generate test bin/cpxctl

cpx-server:
	python3 ./brief/cpx_server.py 8000

deps:
	go mod tidy

generate:
	go generate ./...

test:
	go test ./...
	
bin/cpxctl:
	go build -o $@ main.go	