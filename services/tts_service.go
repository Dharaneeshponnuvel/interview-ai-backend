package services

import (
	"context"
	"encoding/base64"
	"fmt"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// TextToSpeech converts text to audio (MP3 or Base64)
func TextToSpeech(text string, voiceName string, returnType string) (string, []byte, error) {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("TTS client error: %v", err)
	}
	defer client.Close()

	// Build the request
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         voiceName,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	// Call the API
	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("TTS synthesis error: %v", err)
	}

	// Return as Base64 if needed
	if returnType == "base64" {
		audioBase64 := base64.StdEncoding.EncodeToString(resp.AudioContent)
		return audioBase64, nil, nil
	}

	// Otherwise, return raw audio bytes
	return "", resp.AudioContent, nil
}
