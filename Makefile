.PHONY: vendor

include bin/build/make/go.mak

# Run all the specs.
specs: setup-nsq
	go test -race -mod vendor -v -covermode=atomic -coverpkg=./... -coverprofile=test/profile.cov ./...

# Run security checks.
sec:
	gosec -quiet -exclude-dir=test -exclude=G104 ./...

# Setup NSQ.
setup-nsq: delete-nsq create-nsq

# Create NSQ.
create-nsq:
	curl -X POST http://127.0.0.1:4151/topic/create\?topic\=topic
	curl -X POST http://127.0.0.1:4151/channel/create\?topic\=topic\&channel\=channel

# Delete NSQ.
delete-nsq:
	curl -X POST http://127.0.0.1:4151/channel/delete\?topic\=topic\&channel\=channel
	curl -X POST http://127.0.0.1:4151/topic/delete\?topic\=topic

# Send coveralls data.
goveralls: remove-generated-coverage
	goveralls -coverprofile=test/final.cov -service=circle-ci -repotoken=IFpI5rZfnsc2EyZNls8sONCiEB6kFKLiB

# Start the environment.
start:
	bin/build/docker/env start

# Stop the environment.
stop:
	bin/build/docker/env stop
