#pragma once

#include <iostream>

typedef struct s_env {
	std::string sentence_model_path;
}t_env;

bool	tokenize_and_embedd(const std::string& site_content, t_env env);
