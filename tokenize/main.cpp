#include <iostream>
#include <string>
#include <memory>
#include <grpcpp/grpcpp.h>
#include "src/pb/test.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using test::Test;
using test::TestReq;
using test::TestRes;
using test::HashedFile;
using test::Taken;

class TestService final : public Test::Service {
	Status Test(ServerContext* ctx, const TestReq* req, TestRes* res) override {
		std::string msg = req->tm();
		std::cout << "Received message: " << msg << std::endl;
		res->set_tm("Hello from server!");
		(void)ctx; // silence unused parameter warning
		return Status::OK;
	};

	Status TestTokenizeCall(ServerContext* ctx, const HashedFile *req, Taken *res) override {
		std::string hash = req->hash();
		std::cout << "Received hash: " << hash << std::endl;
		// for now just always return true
		// TODO the client prints false for somereason, check why
		res->set_taken(true);
		(void)ctx;
		return Status::OK;
	}
};

void run_server() {
	std::string srv_addr("0.0.0.0:50051");
	TestService service;

	ServerBuilder builder;
	builder.AddListeningPort(srv_addr, grpc::InsecureServerCredentials());
	builder.RegisterService(&service);

	std::unique_ptr<Server> server(builder.BuildAndStart());
	std::cout << "Server listening on " << srv_addr << std::endl;

	//TODO : add signal handling for graceful shutdown (its thread time nooooooo)
	server->Wait();
}

int main() {
	run_server();
	return 0;
}
