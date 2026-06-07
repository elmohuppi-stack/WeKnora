package youtube_transcript

import (
	"context"
	"os"
	"testing"
	"time"
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

	// Test that the provider works even when the original context is cancelled.
	// This simulates the production issue where the HTTP request context gets
	// cancelled (e.g. after 30s) while the Apify call is still in-flight.
	t.Run("FetchTranscriptWithCancelledContext", func(t *testing.T) {
		// Create a context that is already cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		transcript, err := provider.FetchTranscript(ctx, videoID, "en")
		if err != nil {
			t.Fatalf("FetchTranscript with cancelled context failed: %v", err)
		}
		if transcript.Content == "" {
			t.Error("transcript content is empty")
		}
		t.Logf("FetchTranscriptWithCancelledContext - Content (first 200 chars): %s",
			transcript.Content[:min(len(transcript.Content), 200)])
	})

	// Test with a cancelled context for metadata too
	t.Run("FetchMetadataWithCancelledContext", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		meta, err := provider.FetchMetadata(ctx, videoID)
		if err != nil {
			t.Fatalf("FetchMetadata with cancelled context failed: %v", err)
		}
		if meta.Title == "" {
			t.Error("metadata title is empty with cancelled context")
		}
		t.Logf("FetchMetadataWithCancelledContext - Title: %s", meta.Title)
	})

	// Test with a German video to verify language fallback works
	t.Run("FetchTranscriptGermanVideo", func(t *testing.T) {
		// Use a German video - "Was ist Kognitive Kriegsführung?" or similar
		// Using the video the user had issues with
		germanVideoID := "U_pU8pQI3qU"
		transcript, err := provider.FetchTranscript(context.Background(), germanVideoID, "de")
		if err != nil {
			t.Fatalf("FetchTranscript for German video failed: %v", err)
		}
		if transcript.Content == "" {
			t.Error("transcript content is empty for German video")
		}
		t.Logf("German video - Language: %s, Source: %s", transcript.Language, transcript.Source)
		t.Logf("German video - Content (first 200 chars): %s",
			transcript.Content[:min(len(transcript.Content), 200)])
	})

	// Test metadata for the German video
	t.Run("FetchMetadataGermanVideo", func(t *testing.T) {
		germanVideoID := "U_pU8pQI3qU"
		meta, err := provider.FetchMetadata(context.Background(), germanVideoID)
		if err != nil {
			t.Fatalf("FetchMetadata for German video failed: %v", err)
		}
		if meta.Title == "" {
			t.Error("metadata title is empty for German video")
		}
		t.Logf("German video - Title: %s", meta.Title)
		t.Logf("German video - Channel: %s", meta.ChannelName)
		t.Logf("German video - Duration: %ds", meta.Duration)
	})
}

func TestApifyProvider_GermanVideo_FetchTranscript(t *testing.T) {
	apiKey := os.Getenv("APIFY_API_KEY")
	if apiKey == "" {
		t.Skip("APIFY_API_KEY not set, skipping")
	}

	provider, err := NewApifyProvider(apiKey)
	if err != nil {
		t.Fatalf("NewApifyProvider failed: %v", err)
	}

	// The problematic video from the user
	videoID := "ylUtMxtGoLI"

	// Create a context that cancels after 5 seconds to simulate
	// the HTTP request context being cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Logf("Starting FetchMetadata for video %s with 5s timeout...", videoID)
	start := time.Now()

	meta, err := provider.FetchMetadata(ctx, videoID)
	if err != nil {
		// Even if metadata fails (non-fatal), transcript should work
		t.Logf("FetchMetadata failed (non-fatal): %v", err)
	} else {
		t.Logf("Metadata - Title: %s, Channel: %s", meta.Title, meta.ChannelName)
	}

	t.Logf("Starting FetchTranscript for video %s...", videoID)
	transcript, err := provider.FetchTranscript(ctx, videoID, "de")
	elapsed := time.Since(start)
	t.Logf("Elapsed: %v", elapsed)

	if err != nil {
		t.Fatalf("FetchTranscript failed after %v: %v", elapsed, err)
	}
	if transcript.Content == "" {
		t.Error("transcript content is empty")
	}
	t.Logf("Language: %s, Source: %s", transcript.Language, transcript.Source)
	t.Logf("Content length: %d chars", len(transcript.Content))
	t.Logf("Content (first 200 chars): %s",
		transcript.Content[:min(len(transcript.Content), 200)])
}
