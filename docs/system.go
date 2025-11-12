package docs

// System API endpoints documentation
const systemPaths = `
	"/healthcheck": {
		"get": {
			"tags": [
				"system"
			],
			"summary": "API Health Check",
			"description": "Returns the health status of the API",
			"operationId": "healthCheck",
			"responses": {
				"200": {
					"description": "API is healthy",
					"schema": {
						"type": "object",
						"properties": {
							"title": {
								"type": "string",
								"example": "Success"
							},
							"message": {
								"type": "string",
								"example": "API is healthy and running"
							}
						},
						"required": [
							"title",
							"message"
						]
					}
				}
			}
		}
	},
`
