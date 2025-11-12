package docs

// Restaurants API endpoints documentation
const restaurantsPaths = `
	"/restaurants": {
		"get": {
			"tags": [
				"Restaurants"
			],
			"summary": "List all restaurants",
			"description": "Returns a list of restaurants",
			"operationId": "listRestaurants",
			"responses": {
				"200": {
					"description": "Successful operation",
					"schema": {
						"type": "array",
						"items": {
							"$ref": "#/definitions/Restaurant"
						}
					}
				},
				"500": {
					"description": "Internal server error",
					"schema": {
						"$ref": "#/definitions/Error"
					}
				}
			}
		},
		"post": {
			"tags": [
				"Restaurants"
			],
			"summary": "Create a new restaurant",
			"description": "Adds a new restaurant to the system",
			"operationId": "createRestaurant",
			"security": [
				{
					"Bearer": []
				}
			],
			"parameters": [
				{
					"in": "body",
					"name": "restaurant",
					"description": "Restaurant object that needs to be created",
					"required": true,
					"schema": {
						"$ref": "#/definitions/RestaurantInput"
					}
				}
			],
			"responses": {
				"201": {
					"description": "Restaurant created successfully",
					"schema": {
						"$ref": "#/definitions/Restaurant"
					}
				},
				"400": {
					"description": "Invalid input",
					"schema": {
						"$ref": "#/definitions/Error"
					}
				},
				"401": {
					"description": "Unauthorized",
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
