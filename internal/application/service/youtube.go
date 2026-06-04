package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
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
	TranscriptSource string `json:"transcript_source"` // "native" or "generated"
	Transcript       string `json:"transcript"`        // Full transcript text
}

// youTubeURLPatterns matches various YouTube URL formats.
var youTubeURLPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtube\.com\/watch\?(?:.*&)?v=([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtu\.be\/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtube\.com\/embed\/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?youtube\.com\/shorts\/([a-zA-Z0-9_-]{11})`),
}

// supadataBaseURL is the base URL for the Supadata API.
const supadataBaseURL = "https://api.supadata.ai/v1"

// supadataAPIKeyEnv is the environment variable name for the Supadata API key.
const supadataAPIKeyEnv = "SUPADATA_API_KEY"

// getSupadataAPIKey returns the Supadata API key from the environment.
func getSupadataAPIKey() string {
	return os.Getenv(supadataAPIKeyEnv)
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

// supadataHTTPClient is a shared HTTP client for Supadata API calls.
var supadataHTTPClient = &http.Client{Timeout: 60 * time.Second}

// supadataRequest makes an authenticated GET request to the Supadata API.
func supadataRequest(ctx context.Context, endpoint string, query url.Values) (*http.Response, error) {
	apiKey := getSupadataAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("SUPADATA_API_KEY environment variable is not set")
	}

	apiURL := fmt.Sprintf("%s%s?%s", supadataBaseURL, endpoint, query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create Supadata request: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)

	resp, err := supadataHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Supadata request failed: %w", err)
	}
	return resp, nil
}

// --- Supadata API response types ---

// supadataMetadataResponse represents the Supadata Metadata API response.
type supadataMetadataResponse struct {
	Platform    string `json:"platform"`
	Type        string `json:"type"`
	ID          string `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      struct {
		DisplayName string `json:"displayName"`
		AvatarURL   string `json:"avatarUrl"`
	} `json:"author"`
	Stats struct {
		Views    *int64 `json:"views"`
		Likes    *int64 `json:"likes"`
		Comments *int64 `json:"comments"`
		Shares   *int64 `json:"shares"`
	} `json:"stats"`
	Media struct {
		Duration     int    `json:"duration"`
		ThumbnailURL string `json:"thumbnailUrl"`
	} `json:"media"`
	Tags           []string `json:"tags"`
	CreatedAt      string   `json:"createdAt"`
	AdditionalData struct {
		ChannelID string `json:"channelId"`
	} `json:"additionalData"`
}

// supadataTranscriptResponse represents the Supadata Transcript API response (text mode).
type supadataTranscriptResponse struct {
	Content        string   `json:"content"`
	Lang           string   `json:"lang"`
	AvailableLangs []string `json:"availableLangs"`
}

// supadataJobResponse represents a Supadata async job response.
type supadataJobResponse struct {
	JobID string `json:"jobId"`
}

// supadataJobStatusResponse represents the status of an async Supadata job.
type supadataJobStatusResponse struct {
	Status  string `json:"status"`
	Content string `json:"content"`
	Lang    string `json:"lang"`
	Error   string `json:"error,omitempty"`
}

// FetchYouTubeInfo fetches transcript and metadata for a YouTube video using Supadata.
// It uses:
//   - Supadata Metadata API for title, channel, duration, thumbnail, etc.
//   - Supadata Transcript API for captions/transcripts (native or AI-generated)
//
// The API key is read from the SUPADATA_API_KEY environment variable.
// Returns an error if neither transcript nor metadata can be fetched.
func FetchYouTubeInfo(ctx context.Context, videoID string, preferredLang string) (*YouTubeVideoInfo, error) {
	info := &YouTubeVideoInfo{
		VideoID: videoID,
	}

	// Build the canonical YouTube URL
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	// Fetch metadata from Supadata
	metadata, err := fetchSupadataMetadata(ctx, videoURL)
	if err != nil {
		logger.Warnf(ctx, "Failed to fetch Supadata metadata for video %s: %v", videoID, err)
		// Non-fatal: continue with available info
	} else {
		info.Title = metadata.Title
		info.Description = metadata.Description
		info.Duration = metadata.Media.Duration
		info.ThumbnailURL = metadata.Media.ThumbnailURL
		info.ChannelName = metadata.Author.DisplayName
		info.PublishedAt = metadata.CreatedAt
		if metadata.AdditionalData.ChannelID != "" {
			info.ChannelURL = fmt.Sprintf("https://www.youtube.com/channel/%s", metadata.AdditionalData.ChannelID)
		}
	}

	// Fetch transcript from Supadata
	transcriptContent, transcriptLang, transcriptSource, err := fetchSupadataTranscript(ctx, videoURL, preferredLang)
	if err != nil {
		logger.Warnf(ctx, "Failed to fetch transcript for video %s via Supadata: %v", videoID, err)
		return info, fmt.Errorf("failed to fetch transcript: %w", err)
	}

	info.Transcript = transcriptContent
	info.Language = transcriptLang
	info.TranscriptSource = transcriptSource

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

// fetchSupadataMetadata fetches video metadata from the Supadata Metadata API.
func fetchSupadataMetadata(ctx context.Context, videoURL string) (*supadataMetadataResponse, error) {
	query := url.Values{}
	query.Set("url", videoURL)

	resp, err := supadataRequest(ctx, "/metadata", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Supadata metadata API returned status %d: %s", resp.StatusCode, string(body))
	}

	var data supadataMetadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode Supadata metadata response: %w", err)
	}

	return &data, nil
}

// fetchSupadataTranscript fetches transcript from the Supadata Transcript API.
// It first tries with the preferred language, then falls back to any available language.
// Returns (content, language, source, error).
func fetchSupadataTranscript(ctx context.Context, videoURL string, preferredLang string) (string, string, string, error) {
	content, lang, source, err := doFetchSupadataTranscript(ctx, videoURL, preferredLang)
	if err != nil && preferredLang != "" {
		// Fallback: try without language preference
		logger.Infof(ctx, "Transcript fetch with lang=%s failed, retrying without language: %v", preferredLang, err)
		content, lang, source, err = doFetchSupadataTranscript(ctx, videoURL, "")
	}
	if err != nil {
		return "", "", "", err
	}
	return content, lang, source, nil
}

// doFetchSupadataTranscript performs a single transcript fetch request.
func doFetchSupadataTranscript(ctx context.Context, videoURL string, lang string) (string, string, string, error) {
	query := url.Values{}
	query.Set("url", videoURL)
	query.Set("text", "true")
	query.Set("mode", "auto") // auto: try native first, fall back to AI generation
	if lang != "" {
		query.Set("lang", lang)
	}

	resp, err := supadataRequest(ctx, "/transcript", query)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	// Handle async job (HTTP 202) — poll for results
	if resp.StatusCode == http.StatusAccepted {
		var jobResp supadataJobResponse
		if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
			return "", "", "", fmt.Errorf("decode Supadata job response: %w", err)
		}
		if jobResp.JobID == "" {
			return "", "", "", fmt.Errorf("Supadata returned empty job ID")
		}
		logger.Infof(ctx, "Supadata transcript job created: %s", jobResp.JobID)
		return pollSupadataTranscriptJob(ctx, jobResp.JobID)
	}

	// Handle direct response (HTTP 200)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("Supadata transcript API returned status %d: %s", resp.StatusCode, string(body))
	}

	var data supadataTranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", "", fmt.Errorf("decode Supadata transcript response: %w", err)
	}

	if data.Content == "" {
		return "", "", "", fmt.Errorf("Supadata transcript is empty")
	}

	// If we got a direct response with mode=auto, it's a native transcript
	source := "native"
	return data.Content, data.Lang, source, nil
}

// pollSupadataTranscriptJob polls a Supadata transcript job until completion.
func pollSupadataTranscriptJob(ctx context.Context, jobID string) (string, string, string, error) {
	pollURL := fmt.Sprintf("%s/transcript/%s", supadataBaseURL, jobID)
	apiKey := getSupadataAPIKey()

	client := &http.Client{Timeout: 10 * time.Second}
	pollInterval := 1 * time.Second
	maxDuration := 120 * time.Second
	start := time.Now()

	for {
		if time.Since(start) > maxDuration {
			return "", "", "", fmt.Errorf("Supadata transcript job %s timed out after %v", jobID, maxDuration)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, pollURL, nil)
		if err != nil {
			return "", "", "", fmt.Errorf("create poll request: %w", err)
		}
		req.Header.Set("x-api-key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return "", "", "", fmt.Errorf("poll request failed: %w", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			return "", "", "", fmt.Errorf("Supadata transcript job %s not found", jobID)
		}

		var jobStatus supadataJobStatusResponse
		if err := json.NewDecoder(resp.Body).Decode(&jobStatus); err != nil {
			resp.Body.Close()
			return "", "", "", fmt.Errorf("decode job status response: %w", err)
		}
		resp.Body.Close()

		switch jobStatus.Status {
		case "completed":
			if jobStatus.Content == "" {
				return "", "", "", fmt.Errorf("Supadata transcript job %s completed but content is empty", jobID)
			}
			source := "generated"
			return jobStatus.Content, jobStatus.Lang, source, nil
		case "failed":
			errMsg := jobStatus.Error
			if errMsg == "" {
				errMsg = "unknown error"
			}
			return "", "", "", fmt.Errorf("Supadata transcript job %s failed: %s", jobID, errMsg)
		case "queued", "active":
			// Still processing, wait and retry
			time.Sleep(pollInterval)
			continue
		default:
			return "", "", "", fmt.Errorf("unexpected Supadata job status: %s", jobStatus.Status)
		}
	}
}
