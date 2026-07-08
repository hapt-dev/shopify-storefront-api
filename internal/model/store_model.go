package model

import "encoding/json"

// CachedStore is the Redis hash payload written by payment-api store-cache job.
type CachedStore struct {
	Shop             string            `json:"shop"`
	ID               string            `json:"_id"`
	FacebookPixels   []json.RawMessage `json:"facebookPixels"`
	GoogleAnalytics  []json.RawMessage `json:"googleAnalytics"`
	GoogleAds        []json.RawMessage `json:"googleAds"`
	TiktokPixels     []json.RawMessage `json:"tiktokPixels"`
	PinterestTags    []json.RawMessage `json:"pinterestTags"`
	Klaviyo          []json.RawMessage `json:"klaviyo"`
	NewsbreakPixels  []json.RawMessage `json:"newsbreakPixels"`
	SnapchatPixels   []json.RawMessage `json:"snapchatPixels"`
	PixelTrackings   []json.RawMessage `json:"pixelTrackings"`
}

// CachedProduct is the Redis hash payload written by payment-api store-cache job.
type CachedProduct struct {
	ID              string            `json:"_id"`
	Bundle          json.RawMessage   `json:"bundle"`
	FacebookPixels  []json.RawMessage `json:"facebookPixels"`
	GoogleAnalytics []json.RawMessage `json:"googleAnalytics"`
	GoogleAds       []json.RawMessage `json:"googleAds"`
	TiktokPixels    []json.RawMessage `json:"tiktokPixels"`
	PinterestTags   []json.RawMessage `json:"pinterestTags"`
	Klaviyo         []json.RawMessage `json:"klaviyo"`
	NewsbreakPixels []json.RawMessage `json:"newsbreakPixels"`
	SnapchatPixels  []json.RawMessage `json:"snapchatPixels"`
	Pixels          []json.RawMessage `json:"pixels"`
	Reviews         []json.RawMessage `json:"reviews"`
}
