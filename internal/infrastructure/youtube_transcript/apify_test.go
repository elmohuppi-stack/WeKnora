package youtube_transcript

import (
	"context"
	"os"
	"testing"
)

func TestApifyProvider_FetchTranscript(t *testing.T) {
	apiKey := os.Getenv("APIFY_API_KEY")
	if apiKey == "" {
		t.Skip("APIFY_API_KEY not set, skipping")
	}

	provider, err := NewApifyProvider(apiKey)
	if err != nil {
		t.Fatalf("NewApifyProvider failed: %v", err)
	}

	// Use a well-known video with captions
	videoID := "dQw4w9WgXcQ" // Rick Astley - Never Gonna Give You Up

	t.Run("FetchMetadata", func(t *testing.T) {
		meta, err := provider.FetchMetadata(context.Background(), videoID)
		if err != nil {
			t.Fatalf("FetchMetadata failed: %v", err)
		}
		if meta.Title == "" {
			t.Error("metadata title is empty")
		}
		if meta.ChannelName == "" {
			t.Error("metadata channel is empty")
		}
		if meta.Duration <= 0 {
			t.Errorf("metadata duration is %d, expected > 0", meta.Duration)
		}
		t.Logf("Title: %s", meta.Title)
		t.Logf("Channel: %s", meta.ChannelName)
		t.Logf("Duration: %ds", meta.Duration)
		t.Logf("Thumbnail: %s", meta.ThumbnailURL)
	})

	t.Run("FetchTranscript", func(t *testing.T) {
		transcript, err := provider.FetchTranscript(context.Background(), videoID, "en")
		if err != nil {
			t.Fatalf("FetchTranscript failed: %v", err)
		}
		if transcript.Content == "" {
			t.Error("transcript content is empty")
		}
		if transcript.Language == "" {
			t.Error("transcript language is empty")
		}
		t.Logf("Language: %s", transcript.Language)
		t.Logf("Source: %s", transcript.Source)
		t.Logf("Content (first 200 chars): %s", transcript.Content[:min(len(transcript.Content), 200)])
	})
}
