GO_OUT = \033[37;46m
CPP_OUT = \033[37;42m

RESET_COLOR = \033[0m

PROTO_NAME = docuchat.proto

all: go cpp

go:
	@echo "$(GO_OUT)Generating Go code from .proto files...$(RESET_COLOR)"
	protoc	--go_out=../scrapper/pb --go-grpc_out=../scrapper/pb \
			--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
			$(PROTO_NAME)

cpp:
	@echo "$(CPP_OUT)Generating C++ code from .proto files...$(RESET_COLOR)"
	protoc	--cpp_out=../tokenize/src/pb --grpc_out=../tokenize/src/pb \
			--plugin=protoc-gen-grpc=`which grpc_cpp_plugin` \
			$(PROTO_NAME)
