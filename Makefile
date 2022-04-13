.PHONY: vendor tools

help: ## Display this help
	@ echo "Please use \`make <target>' where <target> is one of:"
	@ echo
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-10s\033[0m - %s\n", $$1, $$2}'
	@ echo


setup: dep ## Setup everything

download:
	go mod download

tidy:
	go mod tidy

vendor:
	go mod vendor

get:
	go get $(module)

get-all:
	go get -u all

dep: download tidy vendor ## Setup go deps

lint: ## Lint all the go code
	golangci-lint run --timeout 5m

fix-lint: ## Fix the lint issues in the go code (if possible)
	golangci-lint run --timeout 5m --fix

specs: ## Run all the specs
	go test -mod vendor -v -covermode=atomic -coverpkg=./... -coverprofile=test/profile.cov ./...

remove-generated-coverage:
	cat test/profile.cov | grep -v "test" > test/final.cov

html-coverage: remove-generated-coverage ## Get the HTML coverage for go
	go tool cover -html test/final.cov

func-coverage: remove-generated-coverage  ## Get the func coverage for go
	go tool cover -func test/final.cov

goveralls: remove-generated-coverage ## Send coveralls data
	goveralls -coverprofile=test/final.cov -service=circle-ci -repotoken=IFpI5rZfnsc2EyZNls8sONCiEB6kFKLiB

generate-proto: ## Generate proto for go
	make -C test generate

outdated: ## Check outdated go deps
	go list -u -m -mod=mod -json all | go-mod-outdated -update -direct

update-dep: get tidy vendor ## Update go dep

update-all-deps: get-all tidy vendor ## Update all go dep

setup-nsq: delete-nsq create-nsq ## Setup NSQ

create-nsq: ## Create NSQ
	curl -X POST http://127.0.0.1:4151/topic/create\?topic\=topic
	curl -X POST http://127.0.0.1:4151/channel/create\?topic\=topic\&channel\=channel

delete-nsq: ## Delete NSQ
	curl -X POST http://127.0.0.1:4151/channel/delete\?topic\=topic\&channel\=channel
	curl -X POST http://127.0.0.1:4151/topic/delete\?topic\=topic

start: ## Start the environment
	tools/env start

stop: ## Stop the environment
	tools/env stop
