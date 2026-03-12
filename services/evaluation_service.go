package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type EvaluationResponse struct {
	Score            int      `json:"score"`
	Strengths        []string `json:"strengths"`
	Improvements     []string `json:"improvements"`
	SampleResponse   string   `json:"sample_response"`
	FollowUpQuestion string   `json:"follow_up_question"`
}

func EvaluateAnswer(jobTitle, question, answer string) (EvaluationResponse, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	model := "models/gemini-2.0-flash"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/%s:generateContent?key=%s", model, apiKey)

	prompt := fmt.Sprintf(`
You are an AI interview evaluator.
Evaluate the candidate's answer for the job "%s".

Question: %s
Answer: %s

Return your response *strictly in this JSON format*:

{
  "score": <integer from 0 to 100>,
  "strengths": ["point 1", "point 2"],
  "improvements": ["point 1", "point 2"],
  "sample_response": "ideal answer in 2-3 sentences",
  "follow_up_question": "related follow-up question"
}

Do not include any extra text or explanations outside of JSON.
`, jobTitle, question, answer)

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
		return EvaluationResponse{}, fmt.Errorf("Gemini API error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return EvaluationResponse{}, fmt.Errorf("Gemini returned error: %s", string(body))
	}

	var apiResp map[string]interface{}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return EvaluationResponse{}, fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	text := extractGeminiText(apiResp)

	var result EvaluationResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		fmt.Println("⚠️ Gemini returned non-JSON text:", text)
		return EvaluationResponse{
			Score:            0,
			Strengths:        []string{},
			Improvements:     []string{"AI response not structured."},
			SampleResponse:   text,
			FollowUpQuestion: "Could you clarify your answer?",
		}, nil
	}

	return result, nil
}
