.PHONY: vendor

include bin/build/make/go.mak

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
