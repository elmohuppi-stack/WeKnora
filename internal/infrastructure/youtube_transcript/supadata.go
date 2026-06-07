package youtube_transcript

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

const (
	// supadataBaseURL is the base URL for the Supadata API.
	supadataBaseURL = "https://api.supadata.ai/v1"

	// supadataAPIKeyEnv is the environment variable name for the Supadata API key.
	supadataAPIKeyEnv = "SUPADATA_API_KEY"
)

// supadataHTTPClient is a shared HTTP client for Supadata API calls.
var supadataHTTPClient = &http.Client{Timeout: 60 * time.Second}

// SupadataProvider implements YouTubeTranscriptProvider using the Supadata API.
type SupadataProvider struct {
	apiKey string
}

// NewSupadataProvider creates a new Supadata-based YouTube transcript provider.
func NewSupadataProvider(apiKey string) (interfaces.YouTubeTranscriptProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("SUPADATA_API_KEY is required for Supadata provider")
	}
	return &SupadataProvider{apiKey: apiKey}, nil
}

// Name returns the provider name.
func (p *SupadataProvider) Name() string { return "supadata" }

// FetchMetadata fetches video metadata from the Supadata Metadata API.
func (p *SupadataProvider) FetchMetadata(ctx context.Context, videoID string) (*types.YouTubeMetadataResult, error) {
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	query := url.Values{}
	query.Set("url", videoURL)

	resp, err := p.supadataRequest(ctx, "/metadata", query)
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

	result := &types.YouTubeMetadataResult{
		VideoID:      videoID,
		Title:        data.Title,
		Description:  data.Description,
		Duration:     data.Media.Duration,
		ThumbnailURL: data.Media.ThumbnailURL,
		ChannelName:  data.Author.DisplayName,
		PublishedAt:  data.CreatedAt,
	}
	if data.AdditionalData.ChannelID != "" {
		result.ChannelURL = fmt.Sprintf("https://www.youtube.com/channel/%s", data.AdditionalData.ChannelID)
	}

	return result, nil
}

// FetchTranscript fetches transcript from the Supadata Transcript API.
func (p *SupadataProvider) FetchTranscript(ctx context.Context, videoID string, preferredLang string) (*types.YouTubeTranscriptResult, error) {
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

	content, lang, source, err := p.fetchSupadataTranscript(ctx, videoURL, preferredLang)
	if err != nil {
		return nil, err
	}

	return &types.YouTubeTranscriptResult{
		Content:       content,
		Language:      lang,
		Source:        source,
		IsAIGenerated: source == "generated",
	}, nil
}

// supadataRequest makes an authenticated GET request to the Supadata API.
func (p *SupadataProvider) supadataRequest(ctx context.Context, endpoint string, query url.Values) (*http.Response, error) {
	apiKey := p.apiKey

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

// fetchSupadataTranscript fetches transcript from the Supadata Transcript API.
// It first tries with the preferred language, then falls back to any available language.
func (p *SupadataProvider) fetchSupadataTranscript(ctx context.Context, videoURL string, preferredLang string) (string, string, string, error) {
	content, lang, source, err := p.doSupadataTranscriptFetch(ctx, videoURL, preferredLang)
	if err != nil && preferredLang != "" {
		// Fallback: try without language preference
		logger.Infof(ctx, "Supadata transcript fetch with lang=%s failed, retrying without language: %v", preferredLang, err)
		content, lang, source, err = p.doSupadataTranscriptFetch(ctx, videoURL, "")
	}
	if err != nil {
		return "", "", "", err
	}
	return content, lang, source, nil
}

// doSupadataTranscriptFetch performs a single transcript fetch request.
func (p *SupadataProvider) doSupadataTranscriptFetch(ctx context.Context, videoURL string, lang string) (string, string, string, error) {
	query := url.Values{}
	query.Set("url", videoURL)
	query.Set("text", "true")
	query.Set("mode", "auto") // auto: try native first, fall back to AI generation
	if lang != "" {
		query.Set("lang", lang)
	}

	resp, err := p.supadataRequest(ctx, "/transcript", query)
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
		return p.pollSupadataTranscriptJob(ctx, jobResp.JobID)
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
func (p *SupadataProvider) pollSupadataTranscriptJob(ctx context.Context, jobID string) (string, string, string, error) {
	pollURL := fmt.Sprintf("%s/transcript/%s", supadataBaseURL, jobID)
	apiKey := p.apiKey

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

// Ensure compile-time interface conformance
var _ interfaces.YouTubeTranscriptProvider = (*SupadataProvider)(nil)

// getSupadataAPIKey returns the Supadata API key from the environment.
// Kept for backward compatibility and for use by the old youtube.go code.
func getSupadataAPIKey() string {
	return os.Getenv(supadataAPIKeyEnv)
}
