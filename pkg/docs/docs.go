// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "petsnextdoordev@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/callback/kakao": {
            "get": {
                "description": "Kakao 로그인 콜백을 처리하고, 사용자 기본 정보와 함께 Firebase Custom Token을 발급합니다.",
                "tags": [
                    "auth"
                ],
                "summary": "Kakao 회원가입 콜백 API",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.kakaoCallbackResponse"
                        }
                    }
                }
            }
        },
        "/auth/login/kakao": {
            "get": {
                "tags": [
                    "auth"
                ],
                "summary": "Kakao 로그인 페이지로 redirect 합니다.",
                "responses": {
                    "302": {
                        "description": "Found"
                    }
                }
            }
        },
        "/media/images": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "media"
                ],
                "summary": "이미지를 업로드합니다.",
                "parameters": [
                    {
                        "type": "file",
                        "description": "이미지 파일",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/server.mediaView"
                        }
                    }
                }
            }
        },
        "/media/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "media"
                ],
                "summary": "미디어를 ID로 조회합니다.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "미디어 ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.mediaView"
                        }
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "security": [
                    {
                        "firebase": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "내 프로필 정보를 조회합니다.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.UserResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "firebase": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "내 프로필 정보를 수정합니다.",
                "parameters": [
                    {
                        "description": "프로필 정보 수정 요청",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.UpdateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.UserResponse"
                        }
                    }
                }
            }
        },
        "/users/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "파이어베이스 가입 이후 정보를 입력 받아 유저를 생성합니다.",
                "parameters": [
                    {
                        "description": "사용자 회원가입 요청",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.RegisterUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/server.UserResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.MediaType": {
            "type": "string",
            "enum": [
                "image"
            ],
            "x-enum-varnames": [
                "IMAGE_MEDIA_TYPE"
            ]
        },
        "server.RegisterUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fbProviderType": {
                    "type": "string"
                },
                "fbUid": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                }
            }
        },
        "server.UpdateUserRequest": {
            "type": "object",
            "properties": {
                "nickname": {
                    "type": "string"
                }
            }
        },
        "server.UserResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fbProviderType": {
                    "type": "string"
                },
                "fbUid": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "nickname": {
                    "type": "string"
                }
            }
        },
        "server.kakaoCallbackResponse": {
            "type": "object",
            "properties": {
                "authToken": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fbProviderType": {
                    "type": "string"
                },
                "fbUid": {
                    "type": "string"
                },
                "photoURL": {
                    "type": "string"
                }
            }
        },
        "server.mediaView": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "mediaType": {
                    "$ref": "#/definitions/models.MediaType"
                },
                "url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "이웃집멍냥 API 문서",
	Description:      "이웃집멍냥 백엔드 API 문서입니다.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
