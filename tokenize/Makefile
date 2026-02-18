CXX = g++
# TODO return the -Werror later
CXXFLAGS = -Wall -Wextra -std=c++17 -I.
LDFLAGS = -lgrpc++ -lgrpc -lgpr -lprotobuf -lpthread \
		  -labsl_synchronization -labsl_strings -labsl_base

PROTO_SRCS = src/pb/docuchat.pb.cc src/pb/docuchat.grpc.pb.cc
PROTO_OBJS = $(addprefix $(OBJ_DIR)/, $(notdir $(PROTO_SRCS:.cc=.o)))

SENTENCE_PIECE_LDFLAGS = -lsentencepiece
LDFLAGS += $(SENTENCE_PIECE_LDFLAGS)

ONNXRUNTIME_LDFLAGS = -L/usr/local/onnxruntime/lib -lonnxruntime
LDFLAGS += $(ONNXRUNTIME_LDFLAGS)

ONNXRUNTIME_INCLUDE = -I/usr/local/onnxruntime/include

SRCS = main.cpp tokenize.cpp
OBJS = $(addprefix $(OBJ_DIR)/, $(SRCS:.cpp=.o))

NAME = tokenizer
OBJ_DIR = obj

all: $(OBJ_DIR) $(NAME)

$(OBJ_DIR):
	mkdir -p $(OBJ_DIR)

$(NAME): $(PROTO_OBJS) $(OBJS)
	$(CXX) $(PROTO_OBJS) $(OBJS) $(LDFLAGS) -o $(NAME) -I. $(ONNXRUNTIME_INCLUDE)

$(OBJ_DIR)/%.o: src/pb/%.cc | $(OBJ_DIR)
	$(CXX) $(CXXFLAGS) -c $< -o $@ -I.

$(OBJ_DIR)/%.o: %.cpp | $(OBJ_DIR)
	$(CXX) $(CXXFLAGS) -c $< -o $@ $(ONNXRUNTIME_INCLUDE)

clean:
	rm -rf $(OBJ_DIR)

fclean: clean
	rm -f $(NAME)

re: fclean all

.PHONY: all clean fclean re
