package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/base"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/model"
)

const (
	storeCacheKeyPrefix   = "store:cache"
	productCacheKeyPrefix = "product:cache"
)

var shopifyProductIDPattern = regexp.MustCompile(`/Product/(\d+)`)

type IStoreRepository interface {
	GetStoreCache(ctx context.Context, storeType, shop string) (*model.CachedStore, error)
	GetProductCache(ctx context.Context, shop, productID string) (*model.CachedProduct, error)
}

type storeRepository struct {
	redis *redis.Client
}

func NewStoreRepository(redisClient *redis.Client) IStoreRepository {
	return &storeRepository{redis: redisClient}
}

func StoreCacheKey(storeType, shop string) string {
	return fmt.Sprintf("%s:%s:%s", storeCacheKeyPrefix, storeType, shop)
}

func ProductCacheKey(shop, productID string) string {
	return fmt.Sprintf("%s:%s:%s", productCacheKeyPrefix, shop, NormalizeShopifyProductID(productID))
}

func NormalizeShopifyProductID(productID string) string {
	if productID == "" {
		return productID
	}
	if match := shopifyProductIDPattern.FindStringSubmatch(productID); len(match) == 2 {
		return match[1]
	}
	return productID
}

func (r *storeRepository) GetStoreCache(ctx context.Context, storeType, shop string) (*model.CachedStore, error) {
	key := StoreCacheKey(storeType, shop)
	fields, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis hgetall store cache")
	}
	if len(fields) == 0 {
		return nil, base.NotFound(base.ErrNotFound)
	}

	store := &model.CachedStore{
		Shop:            decodeString(fields["shop"]),
		ID:              decodeString(fields["_id"]),
		FacebookPixels:  decodeArray(fields["facebookPixels"]),
		GoogleAnalytics: decodeArray(fields["googleAnalytics"]),
		GoogleAds:       decodeArray(fields["googleAds"]),
		TiktokPixels:    decodeArray(fields["tiktokPixels"]),
		PinterestTags:   decodeArray(fields["pinterestTags"]),
		Klaviyo:         decodeArray(fields["klaviyo"]),
		NewsbreakPixels: decodeArray(fields["newsbreakPixels"]),
		SnapchatPixels:  decodeArray(fields["snapchatPixels"]),
		PixelTrackings:  decodeArray(fields["pixelTrackings"]),
	}
	return store, nil
}

func (r *storeRepository) GetProductCache(ctx context.Context, shop, productID string) (*model.CachedProduct, error) {
	key := ProductCacheKey(shop, productID)
	fields, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis hgetall product cache")
	}
	if len(fields) == 0 {
		return nil, base.NotFound(base.ErrNotFound)
	}

	product := &model.CachedProduct{
		ID:              decodeString(fields["_id"]),
		Bundle:          decodeRawOrEmptyObject(fields["bundle"]),
		FacebookPixels:  decodeArray(fields["facebookPixels"]),
		GoogleAnalytics: decodeArray(fields["googleAnalytics"]),
		GoogleAds:       decodeArray(fields["googleAds"]),
		TiktokPixels:    decodeArray(fields["tiktokPixels"]),
		PinterestTags:   decodeArray(fields["pinterestTags"]),
		Klaviyo:         decodeArray(fields["klaviyo"]),
		NewsbreakPixels: decodeArray(fields["newsbreakPixels"]),
		SnapchatPixels:  decodeArray(fields["snapchatPixels"]),
		Pixels:          decodeArray(fields["pixels"]),
		Reviews:         decodeArray(fields["reviews"]),
	}
	return product, nil
}

func decodeString(raw string) string {
	if raw == "" {
		return ""
	}
	var s string
	if err := json.Unmarshal([]byte(raw), &s); err == nil {
		return s
	}
	return raw
}

func decodeArray(raw string) []json.RawMessage {
	if raw == "" {
		return []json.RawMessage{}
	}
	var arr []json.RawMessage
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return []json.RawMessage{}
	}
	if arr == nil {
		return []json.RawMessage{}
	}
	return arr
}

func decodeRawOrEmptyObject(raw string) json.RawMessage {
	if raw == "" {
		return json.RawMessage(`{}`)
	}
	var v json.RawMessage
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return json.RawMessage(`{}`)
	}
	if len(v) == 0 || string(v) == "null" {
		return json.RawMessage(`{}`)
	}
	return v
}
