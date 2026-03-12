package services

import (
	"fmt"
	"io"
	"os"
	"strings"

	pdf "github.com/ledongthuc/pdf"
)

// ExtractResumeText reads a PDF and returns plain text
func ExtractResumeText(file io.Reader) (string, error) {
	tmpFile, err := os.CreateTemp("", "resume-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy uploaded file: %v", err)
	}

	tmpFile.Close()

	content, err := readPdf(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read PDF: %v", err)
	}

	text := strings.TrimSpace(content)
	if len(text) < 50 {
		return "", fmt.Errorf("resume text too short or unreadable")
	}

	return text, nil
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf strings.Builder
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	_, err = io.Copy(&buf, b)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
