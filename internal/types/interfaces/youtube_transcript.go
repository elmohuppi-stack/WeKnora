package interfaces

import (
	"context"

	"github.com/Tencent/WeKnora/internal/types"
)

// YouTubeTranscriptProvider defines the interface for YouTube transcript providers.
// Implementations fetch video metadata and transcripts from different backends
// (e.g. Apify, Supadata).
type YouTubeTranscriptProvider interface {
	// Name returns the provider name (e.g. "apify", "supadata").
	Name() string

	// FetchTranscript fetches the transcript for a YouTube video.
	// Returns the transcript content, detected language, and source description.
	FetchTranscript(ctx context.Context, videoID string, preferredLang string) (*types.YouTubeTranscriptResult, error)

	// FetchMetadata fetches video metadata (title, channel, duration, etc.).
	FetchMetadata(ctx context.Context, videoID string) (*types.YouTubeMetadataResult, error)
}
