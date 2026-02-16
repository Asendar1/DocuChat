tokenizer = "tokenize/tokenizer"

all:
	hivemind

run-gateway:
	make -C gateway

run-scrapper:
	make -C scrapper

run-cpp:
	./$(tokenizer)

run-compile-cpp:
	make -C tokenize
	./$(tokenizer)

fclean:
	make -C tokenize fclean

.PHONY: all run-gateway run-scrapper run-compile-cpp
