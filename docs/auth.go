package docs

// Auth API endpoints documentation
const authPaths = `
	"/auth/signup": {
		"post": {
			"tags": [
				"Auth"
			],
			"summary": "User Signup",
			"description": "Creates a new user account",
			"operationId": "signup",
			"parameters": [
				{
					"in": "body",
					"name": "body",
					"description": "Signup request",
					"required": true,
					"schema": {
						"$ref": "#/definitions/SignupInput"
					}
				}
			],
			"responses": {
				"201": {
					"description": "User created successfully",
					"schema": {
						"type": "object",
						"properties": {
							"title": {
								"type": "string",
								"example": "Success"
							},
							"message": {
								"type": "string",
								"example": "User created successfully"
							},
							"user": {
								"type": "object",
								"properties": {
									"id": {
										"type": "string",
										"format": "uuid"
									},
									"name": {
										"type": "string"
									},
									"email": {
										"type": "string",
										"format": "email"
									}
								},
								"required": [
									"id",
									"name",
									"email"
								]
							}
						},
						"required": [
							"title",
							"message",
							"user"
						]
					}
				},
				"400": {
					"description": "Validation error",
					"schema": {
						"$ref": "#/definitions/Error"
					}
				},
				"409": {
					"description": "Duplicate email",
					"schema": {
						"$ref": "#/definitions/Error"
					}
				},
				"500": {
					"description": "Internal server error",
					"schema": {
						"$ref": "#/definitions/Error"
					}
				}
			}
		}
	},
`
