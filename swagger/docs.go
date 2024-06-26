// Package swagger Code generated by swaggo/swag. DO NOT EDIT
package swagger

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
        "/auth/forgot-password": {
            "post": {
                "description": "Validates email if exists in DB, then send an email with verification link to email user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "ForgotPassword",
                "parameters": [
                    {
                        "description": "ForgotPasswordRequest",
                        "name": "ForgotPasswordRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.ForgotPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.StatusResponse"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Validates email and password in request, check if user exists in DB if not throw 404 otherwise compare the request password with hash, then check if user is active, then finds relationships of user with orgs and then generates a JWT token, and returns UserData, Orgs, and Token in response.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "LoginRequest",
                        "name": "LoginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginResponse"
                        }
                    }
                }
            }
        },
        "/auth/reset-password/{token}": {
            "post": {
                "description": "Validates token, new password, and confirm new password, checks if user exists in DB then it updates the password in DB.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "ResetPassword",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "ResetPasswordRequest",
                        "name": "ResetPasswordRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.ResetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.StatusResponse"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Validates email, username, first name, last name, password checks if email exists, if not creates new user and sends email with verification link.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Signup",
                "parameters": [
                    {
                        "description": "SignupRequest",
                        "name": "SignupRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.SignupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.SignupResponse"
                        }
                    }
                }
            }
        },
        "/auth/signup/verify/{token}": {
            "put": {
                "description": "Validates token in param, if token parses valid then user will be verified and be updated in DB.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignupVerify",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.StatusResponse"
                        }
                    }
                }
            }
        },
        "/chapters": {
            "get": {
                "description": "Lists all chapters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chapters"
                ],
                "summary": "Index",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/chapter.ListChapterRequest"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Validates user id and title. If they are up to standard a new chapter will be created. The created chapters ID will be returned.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chapters"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "CreateChapterRequest",
                        "name": "CreateChapterRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/chapter.CreateChapterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chapter.CreateChapterResponse"
                        }
                    }
                }
            }
        },
        "/chapters/{id}": {
            "get": {
                "description": "Fetches a chapter by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chapters"
                ],
                "summary": "Read",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Chapter ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/chapter.ReadChapterRequest"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Updates a chapter",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chapters"
                ],
                "summary": "Update",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Chapter ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "UpdateChapterRequest",
                        "name": "UpdateChapterRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/chapter.UpdateChapterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chapter.UpdateChapterResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a chapter by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chapters"
                ],
                "summary": "Delete",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Chapter ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/chapter.DeleteChapterResponse"
                        }
                    }
                }
            }
        },
        "/oauth/accounts": {
            "get": {
                "description": "Validates user id, will query DB in oauth accounts and returns records.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuthAccount"
                ],
                "summary": "UserOAuthAccounts",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/google.UserAccountResponse"
                            }
                        }
                    }
                }
            }
        },
        "/oauth/google": {
            "get": {
                "description": "Will return the Google OAuth2.0 redirect URL for sign in.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuthAccount"
                ],
                "summary": "GoogleSignIn",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/google.OAuthResponse"
                        }
                    }
                }
            }
        },
        "/oauth/google/callback": {
            "get": {
                "description": "This API is normally automatically called from google redirect of GoogleSignIn endpoint in UI(it requires the code from google to process the process of signing in), and returns the OAuth provider and token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuthAccount"
                ],
                "summary": "GoogleSignInCallback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/google.CallbackResponse"
                        }
                    }
                }
            }
        },
        "/oauth/google/signup": {
            "get": {
                "description": "Will return the Google OAuth2.0 redirect URL for sign up.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuthAccount"
                ],
                "summary": "GoogleSignUp",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/google.OAuthResponse"
                        }
                    }
                }
            }
        },
        "/oauth/google/signup/callback": {
            "get": {
                "description": "This API endpoint is normally automatically called from google redirect of GoogleSignUp endpoint in UI(it requires the code from google to process the process of signing up), and returns the OAuth provider and token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuthAccount"
                ],
                "summary": "GoogleSignUpCallback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/google.CallbackResponse"
                        }
                    }
                }
            }
        },
        "/pages": {
            "get": {
                "description": "Lists all pages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pages"
                ],
                "summary": "Index",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/pages.ListPageRequest"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Validates user id and title. If they are up to standard a new page will be created. The created pages ID will be returned.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pages"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "CreatePageRequest",
                        "name": "CreatePageRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/pages.CreatePageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/pages.CreatePageResponse"
                        }
                    }
                }
            }
        },
        "/pages/{id}": {
            "get": {
                "description": "Fetches a page by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pages"
                ],
                "summary": "Read",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/pages.ReadPageRequest"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Updates a page",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pages"
                ],
                "summary": "Update",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "UpdatePageRequest",
                        "name": "UpdatePageRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/pages.UpdatePageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/pages.UpdatePageResponse"
                        }
                    }
                }
            }
        },
        "/users/update-email": {
            "put": {
                "description": "Validates user id new email then will check if email is the same of exists then if not updates the email of the user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "EmailUpdate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "EmailUpdateRequest",
                        "name": "EmailUpdateRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.EmailUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.StatusResponse"
                        }
                    }
                }
            }
        },
        "/users/update-password": {
            "put": {
                "description": "Validates user id, mode, new password, confirm new password and(or) current password then will set or update password of the user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "PasswordUpdate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization Key(e.g Bearer key)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "PasswordUpdateRequest",
                        "name": "PasswordUpdateRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.PasswordUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.StatusResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.ForgotPasswordRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
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
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "userData": {
                    "$ref": "#/definitions/auth.UserData"
                }
            }
        },
        "auth.ResetPasswordRequest": {
            "type": "object",
            "properties": {
                "confirmNewPassword": {
                    "type": "string"
                },
                "newPassword": {
                    "type": "string"
                }
            }
        },
        "auth.SignupRequest": {
            "type": "object",
            "properties": {
                "confirmPassword": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "auth.SignupResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "auth.StatusResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
                }
            }
        },
        "auth.UserData": {
            "type": "object",
            "properties": {
                "avatarImgUrl": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "profileId": {
                    "type": "integer"
                },
                "role": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "chapter.CreateChapterRequest": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "chapter.CreateChapterResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "chapter.DeleteChapterResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "chapter.ListChapterRequest": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "integer"
                },
                "chapter": {
                    "type": "integer"
                },
                "limit": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "chapter.ReadChapterRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "chapter.UpdateChapterRequest": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "published_at": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "chapter.UpdateChapterResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "google.CallbackResponse": {
            "type": "object",
            "properties": {
                "provider": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "google.OAuthResponse": {
            "type": "object",
            "properties": {
                "redirectUrl": {
                    "type": "string"
                }
            }
        },
        "google.UserAccountResponse": {
            "type": "object",
            "properties": {
                "oAuthAccounts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/google.UserOAuthAccount"
                    }
                },
                "passwordSet": {
                    "type": "boolean"
                }
            }
        },
        "google.UserOAuthAccount": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "provider": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "pages.CreatePageRequest": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "pages.CreatePageResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "pages.ListPageRequest": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "integer"
                },
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "pages.ReadPageRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "pages.UpdatePageRequest": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "published_at": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "pages.UpdatePageResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "users.EmailUpdateRequest": {
            "type": "object",
            "properties": {
                "newEmail": {
                    "type": "string"
                }
            }
        },
        "users.PasswordUpdateRequest": {
            "type": "object",
            "properties": {
                "confirmNewPassword": {
                    "type": "string"
                },
                "currentPassword": {
                    "type": "string"
                },
                "mode": {
                    "type": "string"
                },
                "newPassword": {
                    "type": "string"
                }
            }
        },
        "users.StatusResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
