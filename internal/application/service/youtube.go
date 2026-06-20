package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/infrastructure/youtube_transcript"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// YouTubeVideoInfo holds the metadata and transcript for a YouTube video.
type YouTubeVideoInfo struct {
	VideoID          string `json:"video_id"`
	Title            string `json:"title"`
	ChannelName      string `json:"channel_name"`
	ChannelURL       string `json:"channel_url"`
	Duration         int    `json:"duration"` // seconds
	ThumbnailURL     string `json:"thumbnail_url"`
	Description      string `json:"description"`
	PublishedAt      string `json:"published_at"`
	Language         string `json:"language"`
	TranscriptSource string `json:"transcript_source"` // "native", "auto_generated", "ai_generated"
	Transcript       string `json:"transcript"`        // Full transcript text
}

// youTubeURLPatterns matches various YouTube URL formats.
var youTubeURLPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtube\.com\/watch\?(?:.*&)?v=([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtu\.be\/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtube\.com\/embed\/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtube\.com\/shorts\/([a-zA-Z0-9_-]{11})`),
}

// ExtractYouTubeVideoID extracts the video ID from a YouTube URL.
func ExtractYouTubeVideoID(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Try to extract from "v" parameter for /watch URLs
	if u.Host == "youtube.com" || strings.HasSuffix(u.Host, ".youtube.com") {
		if v := u.Query().Get("v"); v != "" {
			return v, nil
		}
	}

	// Try all patterns
	for _, pattern := range youTubeURLPatterns {
		matches := pattern.FindStringSubmatch(rawURL)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("could not extract video ID from URL: %s", rawURL)
}

// youtubeProvider lazily initializes the YouTube transcript provider.
var youtubeProvider interfaces.YouTubeTranscriptProvider

// getYouTubeProvider returns the configured YouTube transcript provider.
// Initializes on first call using environment variables.
func getYouTubeProvider() (interfaces.YouTubeTranscriptProvider, error) {
	if youtubeProvider != nil {
		return youtubeProvider, nil
	}

	p, err := youtube_transcript.GetProvider()
	if err != nil {
		return nil, err
	}
	youtubeProvider = p
	return youtubeProvider, nil
}

// FetchYouTubeInfo fetches transcript and metadata for a YouTube video.
// It uses the configured provider (Apify by default, Supadata as fallback).
//
// The provider is selected automatically based on environment variables:
//   - APIFY_API_KEY → Apify (primary)
//   - SUPADATA_API_KEY → Supadata (fallback)
//   - YOUTUBE_TRANSCRIPT_PROVIDER → explicit override
//
// Returns an error if neither transcript nor metadata can be fetched.
func FetchYouTubeInfo(ctx context.Context, videoID string, preferredLang string) (*YouTubeVideoInfo, error) {
	prov, err := getYouTubeProvider()
	if err != nil {
		return nil, fmt.Errorf("YouTube transcript provider not configured: %w", err)
	}

	logger.Infof(ctx, "Fetching YouTube info for video %s using provider: %s", videoID, prov.Name())

	info := &YouTubeVideoInfo{
		VideoID: videoID,
	}

	// Fetch metadata (non-fatal: continue if it fails)
	metadata, err := prov.FetchMetadata(ctx, videoID)
	if err != nil {
		logger.Warnf(ctx, "Failed to fetch metadata for video %s via %s: %v", videoID, prov.Name(), err)
		// Fallback: try YouTube oEmbed API (no API key needed, provides title + channel name)
		if oembed, oembedErr := fetchYouTubeOEmbedMetadata(ctx, videoID); oembedErr == nil {
			info.Title = oembed.Title
			info.ChannelName = oembed.ChannelName
			info.ChannelURL = oembed.ChannelURL
			info.ThumbnailURL = oembed.ThumbnailURL
			logger.Infof(ctx, "Fetched metadata via oEmbed fallback for video %s: title=%s channel=%s", videoID, oembed.Title, oembed.ChannelName)
		} else {
			logger.Warnf(ctx, "oEmbed fallback also failed for video %s: %v", videoID, oembedErr)
		}
	} else {
		info.Title = metadata.Title
		info.Description = metadata.Description
		info.Duration = metadata.Duration
		info.ThumbnailURL = metadata.ThumbnailURL
		info.ChannelName = metadata.ChannelName
		info.ChannelURL = metadata.ChannelURL
		info.PublishedAt = metadata.PublishedAt
	}

	// Fetch transcript (fatal: return error if it fails, transcript is the core requirement)
	transcript, err := prov.FetchTranscript(ctx, videoID, preferredLang)
	if err != nil {
		logger.Warnf(ctx, "Failed to fetch transcript for video %s via %s: %v", videoID, prov.Name(), err)
		return info, fmt.Errorf("failed to fetch transcript: %w", err)
	}

	info.Transcript = transcript.Content
	info.Language = transcript.Language
	info.TranscriptSource = transcript.Source

	// Build structured description if metadata was empty
	if info.Description == "" {
		descParts := []string{}
		if info.ChannelName != "" {
			descParts = append(descParts, fmt.Sprintf("Channel: %s", info.ChannelName))
		}
		if info.TranscriptSource != "" {
			descParts = append(descParts, fmt.Sprintf("Transcript: %s", info.TranscriptSource))
		}
		if info.Language != "" {
			descParts = append(descParts, fmt.Sprintf("Language: %s", info.Language))
		}
		info.Description = strings.Join(descParts, " | ")
	}

	return info, nil
}

// oEmbedResponse represents the YouTube oEmbed API response.
type oEmbedResponse struct {
	Title        string `json:"title"`
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Type         string `json:"type"`
}

// fetchYouTubeOEmbedMetadata fetches video metadata from the YouTube oEmbed API.
// This is a lightweight fallback that requires no API key.
func fetchYouTubeOEmbedMetadata(ctx context.Context, videoID string) (*YouTubeVideoInfo, error) {
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	apiURL := fmt.Sprintf("https://www.youtube.com/oembed?url=%s&format=json", url.QueryEscape(videoURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create oEmbed request: %w", err)
	}

	// Use a short timeout since oEmbed is fast
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("oEmbed request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oEmbed returned status %d", resp.StatusCode)
	}

	var oembed oEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&oembed); err != nil {
		return nil, fmt.Errorf("decode oEmbed response: %w", err)
	}

	if oembed.Title == "" {
		return nil, fmt.Errorf("oEmbed returned empty title")
	}

	channelURL := oembed.AuthorURL
	if channelURL == "" {
		channelURL = fmt.Sprintf("https://www.youtube.com/channel/%s", videoID)
	}

	return &YouTubeVideoInfo{
		VideoID:      videoID,
		Title:        oembed.Title,
		ChannelName:  strings.TrimSpace(oembed.AuthorName),
		ChannelURL:   channelURL,
		ThumbnailURL: oembed.ThumbnailURL,
	}, nil
}
