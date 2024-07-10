// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Elmodis",
            "email": "mateusz-gluch@elmodis.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Retrieves and prints API information. This information contains basic API description.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "Retrieves basic API information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/events-summary/data": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data"
                ],
                "summary": "Retrieves script data for asset context",
                "parameters": [
                    {
                        "type": "string",
                        "default": "367,333",
                        "description": "Comma separated list of Asset IDs",
                        "name": "assets",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1708300800:1708387200",
                        "description": "Timestamp range in format {startTs}:{endTs}. Timestamps in seconds, tz=UTC",
                        "name": "ts",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "144h",
                        "description": "Span description string",
                        "name": "span",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "MACHINE,DATA",
                        "description": "Comma separated list of event categories",
                        "name": "category",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "additionalProperties": true
                            }
                        }
                    }
                }
            }
        },
        "/online-summary/data": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data"
                ],
                "summary": "Retrieves online-summary script data for asset context",
                "parameters": [
                    {
                        "type": "string",
                        "default": "367,333",
                        "description": "Comma separated list of Asset IDs",
                        "name": "assets",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "1708300800:1708387200",
                        "description": "Timestamp range in format {startTs}:{endTs}. Timestamps in seconds, tz=UTC",
                        "name": "ts",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "144h",
                        "description": "Span description string",
                        "name": "span",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "additionalProperties": true
                            }
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Retrieves and prints ping message.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "Pings API",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "Retrieves and prints API version. This information contains API app version information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "Retrieves API Version",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "dev-internal-api.elmodis.com",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Internal Scripts API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
