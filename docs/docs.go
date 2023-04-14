// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/integral/generate_key": {
            "post": {
                "tags": [
                    "积分相关接口"
                ],
                "summary": "生成积分充值密钥并发送到指定邮箱",
                "parameters": [
                    {
                        "description": "生成key的请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/integral.BatchGenerateKeyRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/integral/manual/recharge": {
            "post": {
                "tags": [
                    "积分相关接口"
                ],
                "summary": "手动充值",
                "parameters": [
                    {
                        "description": "手动充值请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/integral.ManualRechargeRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/integral/recharge": {
            "post": {
                "tags": [
                    "积分相关接口"
                ],
                "summary": "充值",
                "parameters": [
                    {
                        "description": "充值请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/integral.RechargeRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/integral/record": {
            "post": {
                "tags": [
                    "积分相关接口"
                ],
                "summary": "记录",
                "parameters": [
                    {
                        "description": "请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/integral.RecordRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/openai/v1/audio": {
            "post": {
                "tags": [
                    "代理OpenAI相关接口"
                ],
                "summary": "音频转文字",
                "parameters": [
                    {
                        "type": "file",
                        "description": "音频文件",
                        "name": "audio",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "model",
                        "name": "model",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "language",
                        "name": "language",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "prompt",
                        "name": "prompt",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "temperature",
                        "name": "temperature",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/openai/v1/chat/completions": {
            "post": {
                "tags": [
                    "代理OpenAI相关接口"
                ],
                "summary": "对话",
                "parameters": [
                    {
                        "description": "openai请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/openai.ChatCompletionRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/openai/v1/image": {
            "post": {
                "tags": [
                    "代理OpenAI相关接口"
                ],
                "summary": "生成图片",
                "parameters": [
                    {
                        "description": "openai请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/openai.ImageRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/user/feedback": {
            "post": {
                "tags": [
                    "用户相关接口"
                ],
                "summary": "意见反馈",
                "parameters": [
                    {
                        "description": "反馈信息",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.FeedbackRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/user/login": {
            "post": {
                "tags": [
                    "用户相关接口"
                ],
                "summary": "登录",
                "parameters": [
                    {
                        "description": "登录请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回token",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/profile": {
            "get": {
                "tags": [
                    "用户相关接口"
                ],
                "summary": "获取用户信息",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.ProfileResponse"
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "用户相关接口"
                ],
                "summary": "修改用户信息",
                "parameters": [
                    {
                        "description": "用户信息",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.ProfileRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/v1/user/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户相关接口"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "注册请求参数",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回token",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/verify/send_code": {
            "post": {
                "tags": [
                    "用户相关接口"
                ],
                "summary": "发送验证码",
                "parameters": [
                    {
                        "description": "发送验证码请求参数(默认为email)",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.SendVerificationCodeRequest"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "integral.BatchGenerateKeyRequest": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer",
                    "default": 10
                },
                "type": {
                    "description": "1代表100积分，2代表500积分，3代表1000积分",
                    "type": "integer",
                    "default": 1
                }
            }
        },
        "integral.ManualRechargeRequest": {
            "type": "object",
            "required": [
                "key"
            ],
            "properties": {
                "auth_code": {
                    "description": "允许充值的授权码",
                    "type": "string"
                },
                "key": {
                    "type": "string"
                },
                "to_email": {
                    "type": "string"
                }
            }
        },
        "integral.RechargeRequest": {
            "type": "object",
            "required": [
                "key"
            ],
            "properties": {
                "key": {
                    "type": "string"
                }
            }
        },
        "integral.RecordRequest": {
            "type": "object",
            "required": [
                "model",
                "type"
            ],
            "properties": {
                "model": {
                    "description": "使用模型",
                    "type": "string"
                },
                "size": {
                    "description": "大小，chat为字数，image为尺寸，audio为时长(分钟)",
                    "type": "integer"
                },
                "type": {
                    "description": "计费类型，chat/image/audio",
                    "type": "string"
                }
            }
        },
        "openai.ChatCompletionMessage": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "name": {
                    "description": "This property isn't in the official documentation, but it's in\nthe documentation for the official library for python:\n- https://github.com/openai/openai-python/blob/main/chatml.md\n- https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb",
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "openai.ChatCompletionRequest": {
            "type": "object",
            "properties": {
                "frequency_penalty": {
                    "type": "number"
                },
                "logit_bias": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "max_tokens": {
                    "type": "integer"
                },
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/openai.ChatCompletionMessage"
                    }
                },
                "model": {
                    "type": "string"
                },
                "n": {
                    "type": "integer"
                },
                "presence_penalty": {
                    "type": "number"
                },
                "stop": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "stream": {
                    "type": "boolean"
                },
                "temperature": {
                    "type": "number"
                },
                "top_p": {
                    "type": "number"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "openai.ImageRequest": {
            "type": "object",
            "properties": {
                "n": {
                    "type": "integer"
                },
                "prompt": {
                    "type": "string"
                },
                "response_format": {
                    "type": "string"
                },
                "size": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "user.FeedbackRequest": {
            "type": "object",
            "required": [
                "content",
                "feedbackType"
            ],
            "properties": {
                "content": {
                    "description": "反馈内容",
                    "type": "string"
                },
                "feedbackType": {
                    "description": "反馈类型: 1-问题反馈 2-功能建议 3-咨询 4-其他",
                    "type": "integer"
                }
            }
        },
        "user.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "type": {
                    "description": "可选字段，默认为email",
                    "type": "string",
                    "default": "email"
                }
            }
        },
        "user.ProfileRequest": {
            "type": "object",
            "properties": {
                "avatar": {
                    "description": "头像",
                    "type": "string"
                },
                "nickName": {
                    "description": "昵称",
                    "type": "string"
                }
            }
        },
        "user.ProfileResponse": {
            "type": "object",
            "properties": {
                "avatar": {
                    "description": "头像",
                    "type": "string"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string"
                },
                "integral": {
                    "description": "用户当前积分",
                    "type": "integer"
                },
                "inviteCode": {
                    "description": "邀请码",
                    "type": "string"
                },
                "nickName": {
                    "description": "昵称",
                    "type": "string"
                }
            }
        },
        "user.RegisterRequest": {
            "type": "object",
            "required": [
                "code",
                "email",
                "password"
            ],
            "properties": {
                "code": {
                    "description": "验证码",
                    "type": "string"
                },
                "email": {
                    "description": "邮箱",
                    "type": "string"
                },
                "inviteCode": {
                    "description": "邀请码",
                    "type": "string"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "type": {
                    "description": "可选字段，默认为email",
                    "type": "string",
                    "default": "email"
                }
            }
        },
        "user.SendVerificationCodeRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "type": {
                    "description": "可选字段，默认为email",
                    "type": "string",
                    "default": "email"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
