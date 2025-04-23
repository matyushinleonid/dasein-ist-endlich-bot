package model

var OpenAIResponseSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"days_left": map[string]interface{}{
			"type":        "integer",
			"description": "How many days the user have left in this world",
		},
		"description": map[string]interface{}{
			"type":        "string",
			"description": "How the amount of days left was calculated based on the user's answers",
		},
	},
	"required":             []interface{}{"days_left", "description"},
	"additionalProperties": false,
}

type OpenAIResponse struct {
	DaysLeft    int    `json:"days_left"`
	Description string `json:"description"`
}
