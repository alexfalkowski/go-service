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

# Create certificates:
create-certs:
	mkcert -key-file test/certs/key.pem -cert-file test/certs/cert.pem localhost
	mkcert -client -key-file test/certs/client-key.pem -cert-file test/certs/client-cert.pem localhost
