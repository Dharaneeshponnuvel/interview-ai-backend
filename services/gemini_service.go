package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type GeminiQuestionResponse struct {
	Questions []struct {
		Text string `json:"text"`
	} `json:"questions"`
}

// GenerateQuestions calls Gemini API to create structured interview questions
func GenerateQuestions(jobTitle, jobDescription, resumeText string, num int) (GeminiQuestionResponse, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return GeminiQuestionResponse{}, fmt.Errorf("missing GEMINI_API_KEY environment variable")
	}

	model := "models/gemini-2.0-flash"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/%s:generateContent?key=%s", model, apiKey)

	prompt := fmt.Sprintf(`
You are an AI interviewer.

Create exactly %d highly relevant interview questions for the job title "%s".

Consider the following:
- Job Description: "%s"
- Candidate Resume Summary: "%s"

Each question should test the candidate's understanding, experience, and skills relevant to this specific job.

Return strictly in this JSON format (no markdown or commentary):

{
  "questions": [
    {"text": "Question 1"},
    {"text": "Question 2"},
    {"text": "Question 3"}
  ]
}
`, num, jobTitle, jobDescription, resumeText)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return GeminiQuestionResponse{}, fmt.Errorf("Gemini API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return GeminiQuestionResponse{}, fmt.Errorf("Gemini error: %s", string(body))
	}

	var apiResp map[string]interface{}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return GeminiQuestionResponse{}, fmt.Errorf("Gemini response parse error: %v", err)
	}

	text := extractGeminiText(apiResp)
	fmt.Println("🧠 Cleaned Gemini text:\n", text)

	var result GeminiQuestionResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		fmt.Println("⚠️ Gemini returned invalid JSON. Returning fallback.")
		return GeminiQuestionResponse{
			Questions: []struct {
				Text string `json:"text"`
			}{
				{Text: text},
			},
		}, nil
	}

	return result, nil
}

// Extracts clean text from Gemini’s API response
func extractGeminiText(apiResp map[string]interface{}) string {
	if c, ok := apiResp["candidates"].([]interface{}); ok && len(c) > 0 {
		candidate := c[0].(map[string]interface{})
		if content, ok := candidate["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				if part, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := part["text"].(string); ok {
						return cleanGeminiJSON(text)
					}
				}
			}
		}
	}
	return ""
}

// Removes ```json fences, markdown, and unnecessary prefixes
func cleanGeminiJSON(raw string) string {
	cleaned := raw
	re := regexp.MustCompile("(?s)```(json)?(.*?)```")
	cleaned = re.ReplaceAllString(cleaned, "$2")
	cleaned = strings.TrimPrefix(cleaned, "json")
	cleaned = strings.TrimPrefix(cleaned, "JSON")
	cleaned = strings.ReplaceAll(cleaned, "```", "")
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}
