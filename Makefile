BINNAME ?= icbc-checker


.PHONY: docker-build
docker-build: build test
	docker build -f ./build/Dockerfile -t icbc-checker .


.PHONY: build
build:
	mkdir -p ./dist
	@echo "Building icbc-checker binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ./dist/$(BINNAME) ./cmd
	@echo "Done!"


.PHONY: test
test:
	@echo "Running unit tests..."
	go test -v ./...