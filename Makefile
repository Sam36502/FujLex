BIN = "bin/fujlex"
IMAGE = "clexicon_web"
CONTAINER = "lex"
REPO = "bismarck6502/"${IMAGE}
VERSION = "1.0.0"

.PHONY: help gen build image up down publish
help: ## Lists all commands and their descriptions
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: gen ## Builds the executable for linux
	@echo "### Building Binary... ###"
	@go build -o ${BIN} ./src/main.go

run: build ## Runs the API locally
	@echo "### Starting API... ###"
	@./${BIN}
	
image: build ## Builds docker image
	@echo "### Building Docker Image... ###"
	@docker build --no-cache -t ${IMAGE} .

up: image ## Start docker image locally
	@echo "### Starting Local Container... ###"
	@docker run -d \
		-p	1919:1919 \
	       	--name ${CONTAINER} \
		-v '/home/bismarck/fujlex/static:/static' \
		-v '/home/bismarck/fujlex/tmpl:/tmpl' \
		${IMAGE}:latest

down: ## Stop local docker image
	@echo "### Stopping Local Container... ###"
	@docker stop ${CONTAINER}
	@docker rm ${CONTAINER}

publish: image ## Publishes docker image to dockerhub
	@echo "### Publishing Docker Image... ###"
	@docker tag ${IMAGE} ${REPO}:${VERSION}
	@docker push ${REPO}:${VERSION}
