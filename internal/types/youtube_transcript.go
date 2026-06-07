package types

// YouTubeTranscriptProviderType identifies a YouTube transcript provider backend.
type YouTubeTranscriptProviderType string

const (
	YouTubeTranscriptProviderTypeApify    YouTubeTranscriptProviderType = "apify"
	YouTubeTranscriptProviderTypeSupadata YouTubeTranscriptProviderType = "supadata"
)

// YouTubeTranscriptResult holds the transcript fetched from a YouTube video.
type YouTubeTranscriptResult struct {
	Content       string `json:"content"`
	Language      string `json:"language"`
	Source        string `json:"source"`          // "native", "auto_generated", "ai_generated"
	IsAIGenerated bool   `json:"is_ai_generated"` // true if AI-transcribed (no native captions)
}

// YouTubeMetadataResult holds video metadata fetched from a provider.
type YouTubeMetadataResult struct {
	VideoID      string   `json:"video_id"`
	Title        string   `json:"title"`
	ChannelName  string   `json:"channel_name"`
	ChannelURL   string   `json:"channel_url"`
	Duration     int      `json:"duration"` // seconds
	ThumbnailURL string   `json:"thumbnail_url"`
	Description  string   `json:"description"`
	PublishedAt  string   `json:"published_at"`
	Tags         []string `json:"tags,omitempty"`
}
