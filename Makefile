NAME=goophermart
VERSION=0.1


go-compile: clean build-server

.PHONY: build-server
build-server:
	@go build  -o server ./cmd/gophermart/main.go

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f ./server

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: help
all: help

# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
