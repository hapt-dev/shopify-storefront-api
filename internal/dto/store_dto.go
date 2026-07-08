package dto

import "encoding/json"

const DefaultStoreType = "SHOPIFY"

type StoreQuery struct {
	Shop      string `query:"shop" validate:"required"`
	ProductID string `query:"productId"`
	StoreType string `query:"storeType"`
}

type ProductQuery struct {
	Shop      string `query:"shop" validate:"required"`
	ProductID string `query:"productId" validate:"required"`
}

type InsightsResponse struct {
	Pixels  []json.RawMessage `json:"pixels"`
	Reviews []json.RawMessage `json:"reviews"`
}

type PixelsResponse struct {
	Shop            string            `json:"shop"`
	ID              string            `json:"_id"`
	FacebookPixels  []json.RawMessage `json:"facebookPixels"`
	GoogleAnalytics []json.RawMessage `json:"googleAnalytics"`
	GoogleAds       []json.RawMessage `json:"googleAds"`
	TiktokPixels    []json.RawMessage `json:"tiktokPixels"`
	PinterestTags   []json.RawMessage `json:"pinterestTags"`
	Klaviyo         []json.RawMessage `json:"klaviyo"`
	NewsbreakPixels []json.RawMessage `json:"newsbreakPixels"`
	SnapchatPixels  []json.RawMessage `json:"snapchatPixels"`
}

type ProductResponse struct {
	ID     string          `json:"_id"`
	Bundle json.RawMessage `json:"bundle"`
}
