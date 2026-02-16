#include <iostream>
#include <string>

#include <sentencepiece_processor.h>
#include <vector>

#include "tokenize.hpp"

class Tokenizer {
	private:
		std::vector<std::string> 				_tokens;
		sentencepiece::SentencePieceProcessor	_processor;

	public:
		// This constructor throws a runtime_error if the model fails to load
		Tokenizer(const std::string &model_path) {
			const auto status = _processor.Load(model_path);
			if (!status.ok()) {
				std::cerr << "Failed to load the model from " << model_path << ": " << status.ToString() << std::endl;
				throw std::runtime_error("Model loading failed");
			}
		}

		std::vector<std::string> tokenize(const std::string &text) {
			std::vector<std::string> tokens;
			const auto status = _processor.Encode(text, &tokens);
			if (!status.ok()) {
				std::cerr << "Tokenization failed for text: " << text << ": " << status.ToString() << std::endl;
				return {};
			}
			// test
			std::vector<std::string>::iterator it = tokens.begin();
			int count = 0;
			while (it != tokens.end() && count < 100) {
				std::cout << *it << std::endl;
				++it;
				++count;
			}
			return tokens;
		}


};


bool	tokenize_and_embedd(const std::string &site_content, t_env env) {
	try {
		Tokenizer tokenizer(env.sentence_model_path);
		tokenizer.tokenize(site_content);
	}
	catch (const std::exception &e) {
		std::cerr << "Error in tokenize_and_embedd: " << e.what() << std::endl;
		return false;
	}
	return true;
}
