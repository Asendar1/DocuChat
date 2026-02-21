#include <iostream>
#include <string>
#include <memory>
#include <grpcpp/grpcpp.h>
#include <csignal>
#include <thread>

#include "src/pb/docuchat.grpc.pb.h"
#include "tokenize.hpp"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using docuchat::DocumentProcessor;
using docuchat::VectorSearch;
using docuchat::FileRequest;
using docuchat::HashedFile;
using docuchat::SearchRequest;
using docuchat::SearchResponse;
using docuchat::ProcessResponse;
using docuchat::TestReq;
using docuchat::TestRes;

t_env env;

// global server isntance for signal handling
std::unique_ptr<Server> g_server;
std::atomic<bool> g_shutdown_requested(false);

void signal_handler(int signal)
{
	g_shutdown_requested.store(true);
}

class DocumentProcessorServiceImpl final : public DocumentProcessor::Service {
	grpc::Status Test(ServerContext* ctx, const TestReq* req, TestRes* res) override {
		std::cout << "Received test request with message: " << req->tm() << std::endl;
		res->set_tm("Hello from the server!");
		return grpc::Status::OK;
	}

	grpc::Status ProcessFile(ServerContext* ctx, const FileRequest* req, ProcessResponse* res) override {
		std::cout << "Received file hash: " << req->hash() << std::endl;
		std::string content(req->content().substr(0, 100));
		std::cout << "Received file content (first 100 chars): " << content << std::endl;
		// TODO Check if the file were already processed, for now i'll just return no
		if (!tokenize_and_embedd(content, env)) {
			return grpc::Status(grpc::StatusCode::INTERNAL, "Failed to process file");
		}
		res->set_success(true);
		res->set_already_exists(true);
		res->set_message("File was already processed");
		res->set_tokens(42);
		return grpc::Status::OK;
	}
};

class VectorSearchServiceImpl final : public VectorSearch::Service {
	grpc::Status Search(ServerContext* ctx, const SearchRequest* req, SearchResponse* res) override {
		std::cout << "Received search request with query: " << req->query() << std::endl;
		// TODO : implement search logic, for now i'll just return an empty response
		return grpc::Status::OK;
	}

	grpc::Status GetDocument(ServerContext* ctx, const HashedFile* req, docuchat::Document* res) override {
		std::cout << "Received document request with hash: " << req->hash() << std::endl;
		// TODO : implement get document logic, for now i'll just return an empty response
		return grpc::Status::OK;
	}

	grpc::Status DeleteDocument(ServerContext* ctx, const HashedFile* req, docuchat::DeleteResponse* res) override {
		std::cout << "Received delete request with hash: " << req->hash() << std::endl;
		// TODO : implement delete document logic, for now i'll just return an empty response
		return grpc::Status::OK;
	}
};

void run_server() {
	std::string srv_addr("0.0.0.0:50051");
	DocumentProcessorServiceImpl doc_service;
	VectorSearchServiceImpl search_service;

	ServerBuilder builder;
	builder.AddListeningPort(srv_addr, grpc::InsecureServerCredentials());
	builder.RegisterService(&doc_service);
	builder.RegisterService(&search_service);

	g_server = builder.BuildAndStart();
	std::cout << "Server listening on " << srv_addr << std::endl;

	//TODO : add signal handling for graceful shutdown (its thread time nooooooo)
	std::signal(SIGINT, signal_handler);
	std::signal(SIGTERM, signal_handler);

	std::thread monitor_thread([&]() {
		while(!g_shutdown_requested.load())
		{
			std::this_thread::sleep_for(std::chrono::seconds(1));
		}
		if (g_server)
		{
			std::cout << "Shutdown signal received, shutting down server..." << std::endl;
			g_server->Shutdown();
		}
	});

	g_server->Wait();
	monitor_thread.join();
	std::cout << "Server shutdown complete" << std::endl;
}

int main() {
	env.sentence_model_path = std::getenv("SENTENCE_MODEL_PATH") ? std::getenv("SENTENCE_MODEL_PATH") : "model/spiece.model";

	run_server();
	return 0;
}
