package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"google.golang.org/grpc/status"
)

func SpeechToText(audioData []byte, filename string) (string, error) {
	fmt.Println("🎤 [DEBUG] SpeechToText called with file:", filename)

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("STT client error: %v", err)
	}
	defer client.Close()

	ext := strings.ToLower(filepath.Ext(filename))
	fmt.Println("🎤 [DEBUG] Detected file extension:", ext)

	var encoding speechpb.RecognitionConfig_AudioEncoding
	switch ext {
	case ".wav":
		encoding = speechpb.RecognitionConfig_LINEAR16
	case ".flac":
		encoding = speechpb.RecognitionConfig_FLAC
	case ".mp3", ".m4a":
		encoding = speechpb.RecognitionConfig_ENCODING_UNSPECIFIED
	default:
		encoding = speechpb.RecognitionConfig_ENCODING_UNSPECIFIED
	}

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     encoding,
			LanguageCode: "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioData},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			fmt.Println("❌ [Google STT Error]:", st.Message())
		} else {
			fmt.Println("❌ [STT Recognition Error]:", err)
		}
		return "", fmt.Errorf("STT recognition error: %v", err)
	}

	transcript := ""
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			transcript += alt.Transcript + " "
		}
	}

	if transcript == "" {
		return "", fmt.Errorf("no transcript generated — check audio quality or encoding")
	}

	fmt.Println("✅ [DEBUG] Transcript:", transcript)
	return strings.TrimSpace(transcript), nil
}
