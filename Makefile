test = "tokenize/test"

all:
	hivemind

run-gateway:
	make -C gateway

run-scrapper:
	make -C scrapper

run-compile-cpp:
	make -C tokenize
	./$(test)

.PHONY: all run-gateway run-scrapper run-compile-cpp
