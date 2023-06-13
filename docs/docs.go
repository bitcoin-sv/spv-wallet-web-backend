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
        "/api/v1/sign-in": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Sign in user",
                "parameters": [
                    {
                        "description": "User sign in data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/access.SignInUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/access.SignInResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/sign-out": {
            "post": {
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Sign out user",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/transaction": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transaction"
                ],
                "summary": "Get all transactions.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/buxclient.Transaction"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/transaction/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transaction"
                ],
                "summary": "Get transaction by id.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Transaction id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/buxclient.FullTransaction"
                        }
                    }
                }
            }
        },
        "/api/v1/user": {
            "post": {
                "description": "Register new user with given data, paymail is created based on username from sended email.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.RegisterUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.RegisterResponse"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "status"
                ],
                "summary": "Check the status of the server",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/user": {
            "get": {
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.UserResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "access.SignInResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "$ref": "#/definitions/users.Balance"
                },
                "paymail": {
                    "type": "string"
                }
            }
        },
        "access.SignInUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "buxclient.FullTransaction": {
            "type": "object",
            "properties": {
                "blockHash": {
                    "type": "string"
                },
                "blockHeight": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "direction": {
                    "type": "string"
                },
                "fee": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "numberOfInputs": {
                    "type": "integer"
                },
                "numberOfOutputs": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "totalValue": {
                    "type": "integer"
                }
            }
        },
        "buxclient.Transaction": {
            "type": "object",
            "properties": {
                "direction": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "totalValue": {
                    "type": "integer"
                }
            }
        },
        "users.Balance": {
            "type": "object",
            "properties": {
                "bsv": {
                    "type": "number"
                },
                "satoshis": {
                    "type": "integer"
                },
                "usd": {
                    "type": "number"
                }
            }
        },
        "users.RegisterResponse": {
            "type": "object",
            "properties": {
                "mnemonic": {
                    "type": "string"
                },
                "paymail": {
                    "type": "string"
                }
            }
        },
        "users.RegisterUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "passwordConfirmation": {
                    "type": "string"
                }
            }
        },
        "users.UserResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "$ref": "#/definitions/users.Balance"
                },
                "email": {
                    "type": "string"
                },
                "paymail": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Bux Wallet API",
	Description:      "This is an API for bux wallet.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
