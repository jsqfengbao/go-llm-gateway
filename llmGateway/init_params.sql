-- 添加 LLM 抽象层相关的系统参数

-- 检查并添加 use_llm_abstraction 参数（是否使用 LLM 抽象层）
INSERT INTO sys_params (created_at, updated_at, name, `key`, `value`, `desc`)
SELECT NOW(), NOW(), '是否使用LLM抽象层', 'use_llm_abstraction', 'true', '是否使用LLM抽象层，true或false'
WHERE NOT EXISTS (
    SELECT 1 FROM sys_params WHERE `key` = 'use_llm_abstraction'
);

-- 检查并添加 default_llm_type 参数（默认 LLM 类型）
INSERT INTO sys_params (created_at, updated_at, name, `key`, `value`, `desc`)
SELECT NOW(), NOW(), '默认LLM类型', 'default_llm_type', 'doubao', '默认使用的LLM类型，可选值：deepseek、openai、doubao'
WHERE NOT EXISTS (
    SELECT 1 FROM sys_params WHERE `key` = 'default_llm_type'
);

-- 检查并添加 openai_apiKey 参数（OpenAI API 密钥）
INSERT INTO sys_params (created_at, updated_at, name, `key`, `value`, `desc`)
SELECT NOW(), NOW(), 'OpenAI API密钥', 'openai_apiKey', '', 'OpenAI API密钥'
WHERE NOT EXISTS (
    SELECT 1 FROM sys_params WHERE `key` = 'openai_apiKey'
);

-- 检查并添加 openai_model 参数（OpenAI 模型名称）
INSERT INTO sys_params (created_at, updated_at, name, `key`, `value`, `desc`)
SELECT NOW(), NOW(), 'OpenAI模型名称', 'openai_model', 'gpt-4-turbo', 'OpenAI模型名称'
WHERE NOT EXISTS (
    SELECT 1 FROM sys_params WHERE `key` = 'openai_model'
);

-- 检查并添加 doubao_apiKey 参数（豆包 API 密钥）
INSERT INTO sys_params (created_at, updated_at, name, `key`, `value`, `desc`)
SELECT NOW(), NOW(), '豆包API密钥', 'doubao_apiKey', '', '豆包API密钥'
WHERE NOT EXISTS (
    SELECT 1 FROM sys_params WHERE `key` = 'doubao_apiKey'
);

-- 检查并添加 doubao_model 参数（豆包模型名称）
INSERT INTO sys_params (created_at, updated_at, name, `key`, `value`, `desc`)
SELECT NOW(), NOW(), '豆包模型名称', 'doubao_model', 'doubao-lite', '豆包模型名称'
WHERE NOT EXISTS (
    SELECT 1 FROM sys_params WHERE `key` = 'doubao_model'
);

-- 检查并添加 timeOut 参数（超时时间）
INSERT INTO sys_params (created_at, updated_at, name, `key`, `value`, `desc`)
SELECT NOW(), NOW(), '超时时间', 'timeOut', '30', '请求超时时间（秒）'
WHERE NOT EXISTS (
    SELECT 1 FROM sys_params WHERE `key` = 'timeOut'
);

-- 更新参数说明
UPDATE sys_params
SET `desc` = 'DeepSeek API密钥'
WHERE `key` = 'deepseek_apiKey'
  AND `desc` != 'DeepSeek API密钥';
