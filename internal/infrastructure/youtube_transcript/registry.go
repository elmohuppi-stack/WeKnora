package youtube_transcript

import (
	"fmt"
	"os"
	"sync"

	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

// ProviderFactory creates a new YouTube transcript provider instance.
type ProviderFactory func(apiKey string) (interfaces.YouTubeTranscriptProvider, error)

// Registry manages YouTube transcript provider type registrations.
type Registry struct {
	factories map[string]ProviderFactory
	mu        sync.RWMutex
}

// NewRegistry creates a new YouTube transcript provider registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]ProviderFactory),
	}
}

// Register registers a provider type factory by ID.
func (r *Registry) Register(id string, factory ProviderFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[id] = factory
}

// CreateProvider creates a provider instance by type with the given API key.
func (r *Registry) CreateProvider(providerType string, apiKey string) (interfaces.YouTubeTranscriptProvider, error) {
	r.mu.RLock()
	factory, ok := r.factories[providerType]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("YouTube transcript provider type %s not registered", providerType)
	}
	return factory(apiKey)
}

// DefaultProviderEnv is the optional env var to force a specific provider.
const DefaultProviderEnv = "YOUTUBE_TRANSCRIPT_PROVIDER"

// GetDefaultProvider returns the default YouTube transcript provider based on environment variables.
// Priority:
//  1. YOUTUBE_TRANSCRIPT_PROVIDER env var (explicit choice)
//  2. APIFY_API_KEY env var → apify
//  3. SUPADATA_API_KEY env var → supadata
//  4. None set → error
func GetDefaultProvider(r *Registry) (interfaces.YouTubeTranscriptProvider, error) {
	// Explicit provider selection via env var
	if forced := os.Getenv(DefaultProviderEnv); forced != "" {
		apiKey := os.Getenv(apiKeyEnvForProvider(forced))
		if apiKey == "" {
			return nil, fmt.Errorf("YOUTUBE_TRANSCRIPT_PROVIDER is set to %q but %s is not set",
				forced, apiKeyEnvForProvider(forced))
		}
		return r.CreateProvider(forced, apiKey)
	}

	// Auto-detect: try Apify first, then Supadata
	if apiKey := os.Getenv("APIFY_API_KEY"); apiKey != "" {
		return r.CreateProvider("apify", apiKey)
	}
	if apiKey := os.Getenv("SUPADATA_API_KEY"); apiKey != "" {
		return r.CreateProvider("supadata", apiKey)
	}

	return nil, fmt.Errorf("no YouTube transcript provider configured: set APIFY_API_KEY or SUPADATA_API_KEY in .env")
}

// apiKeyEnvForProvider returns the env var name for a provider's API key.
func apiKeyEnvForProvider(providerType string) string {
	switch providerType {
	case "apify":
		return "APIFY_API_KEY"
	case "supadata":
		return "SUPADATA_API_KEY"
	default:
		return providerType + "_API_KEY"
	}
}

// defaultRegistry is the package-level registry initialized once.
var defaultRegistry *Registry
var registryOnce sync.Once

// getDefaultRegistry returns the singleton registry with default providers registered.
func getDefaultRegistry() *Registry {
	registryOnce.Do(func() {
		defaultRegistry = NewRegistry()
		defaultRegistry.Register("supadata", NewSupadataProvider)
		defaultRegistry.Register("apify", NewApifyProvider)
	})
	return defaultRegistry
}

// GetProvider returns the default provider (cached). Convenience function for use in youtube.go.
func GetProvider() (interfaces.YouTubeTranscriptProvider, error) {
	return GetDefaultProvider(getDefaultRegistry())
}
