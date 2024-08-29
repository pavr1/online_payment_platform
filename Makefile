# Top-level Makefile

all: auth bank payment_platform

auth:
	$(MAKE) -C auth build

bank:
	$(MAKE) -C bank build

payment_platform:
	$(MAKE) -C payment_platform build

.PHONY: all auth bank payment_platform