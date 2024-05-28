include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams: crypto-diagram telemetry-diagram transport-diagram

crypto-diagram:
	$(MAKE) package=crypto create-diagram

telemetry-diagram:
	$(MAKE) package=telemetry create-diagram

transport-diagram:
	$(MAKE) package=transport create-diagram

create-diagram:
	goda graph github.com/alexfalkowski/go-service/$(package)/... | dot -Tpng -o assets/$(package).png
