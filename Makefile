.PHONY: init-env build dev build-base tests

init-env:
	./scripts/init_env.sh

build: init-env
	docker-compose build

dev: build
	docker-compose up -d

build-base:
	docker build . -f build/Dockerfile --target builder -t myticket_base:latest

tests: build-base
	docker run -t myticket_base:latest go test -cover ./...
