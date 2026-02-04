all:
	hivemind

run-gateway:
	make -C gateway

run-scrapper:
	make -C scrapper

.PHONY: all run-gateway run-scrapper
