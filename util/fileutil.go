package util

import (
	"log"
	"os"
)

func LoadJS(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read js file: %v", err)
	}
	return string(content)
}
