package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/base"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/dto"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/repository"
)

type IStoreService interface {
	GetInsights(ctx context.Context, q dto.StoreQuery) (*dto.InsightsResponse, error)
	GetPixels(ctx context.Context, q dto.StoreQuery) (*dto.PixelsResponse, error)
	GetProduct(ctx context.Context, q dto.ProductQuery) (*dto.ProductResponse, error)
	GetReviews(ctx context.Context, q dto.ProductQuery) ([]json.RawMessage, error)
}

type storeService struct {
	repo repository.IStoreRepository
}

func NewStoreService(repo repository.IStoreRepository) IStoreService {
	return &storeService{repo: repo}
}

func (s *storeService) GetInsights(ctx context.Context, q dto.StoreQuery) (*dto.InsightsResponse, error) {
	storeType := normalizeStoreType(q.StoreType)

	if q.ProductID == "" {
		store, err := s.repo.GetStoreCache(ctx, storeType, q.Shop)
		if err != nil {
			return nil, err
		}
		return &dto.InsightsResponse{
			Pixels:  emptyIfNil(store.PixelTrackings),
			Reviews: []json.RawMessage{},
		}, nil
	}

	product, err := s.repo.GetProductCache(ctx, q.Shop, q.ProductID)
	if err != nil {
		return nil, err
	}
	return &dto.InsightsResponse{
		Pixels:  emptyIfNil(product.Pixels),
		Reviews: emptyIfNil(product.Reviews),
	}, nil
}

func (s *storeService) GetPixels(ctx context.Context, q dto.StoreQuery) (*dto.PixelsResponse, error) {
	storeType := normalizeStoreType(q.StoreType)

	store, err := s.repo.GetStoreCache(ctx, storeType, q.Shop)
	if err != nil {
		return nil, err
	}

	result := &dto.PixelsResponse{
		Shop:            store.Shop,
		ID:              store.ID,
		FacebookPixels:  emptyIfNil(store.FacebookPixels),
		GoogleAnalytics: emptyIfNil(store.GoogleAnalytics),
		GoogleAds:       emptyIfNil(store.GoogleAds),
		TiktokPixels:    emptyIfNil(store.TiktokPixels),
		PinterestTags:   emptyIfNil(store.PinterestTags),
		Klaviyo:         emptyIfNil(store.Klaviyo),
		NewsbreakPixels: emptyIfNil(store.NewsbreakPixels),
		SnapchatPixels:  emptyIfNil(store.SnapchatPixels),
	}

	if q.ProductID == "" {
		return result, nil
	}

	product, err := s.repo.GetProductCache(ctx, q.Shop, q.ProductID)
	if err != nil {
		if errors.Is(err, base.ErrNotFound) {
			return result, nil
		}
		return nil, err
	}

	result.FacebookPixels = append(result.FacebookPixels, emptyIfNil(product.FacebookPixels)...)
	result.GoogleAnalytics = append(result.GoogleAnalytics, emptyIfNil(product.GoogleAnalytics)...)
	result.GoogleAds = append(result.GoogleAds, emptyIfNil(product.GoogleAds)...)
	result.TiktokPixels = append(result.TiktokPixels, emptyIfNil(product.TiktokPixels)...)
	result.PinterestTags = append(result.PinterestTags, emptyIfNil(product.PinterestTags)...)
	result.Klaviyo = append(result.Klaviyo, emptyIfNil(product.Klaviyo)...)
	result.NewsbreakPixels = append(result.NewsbreakPixels, emptyIfNil(product.NewsbreakPixels)...)
	result.SnapchatPixels = append(result.SnapchatPixels, emptyIfNil(product.SnapchatPixels)...)

	return result, nil
}

func (s *storeService) GetProduct(ctx context.Context, q dto.ProductQuery) (*dto.ProductResponse, error) {
	product, err := s.repo.GetProductCache(ctx, q.Shop, q.ProductID)
	if err != nil {
		return nil, err
	}
	bundle := product.Bundle
	if len(bundle) == 0 {
		bundle = json.RawMessage(`{}`)
	}
	return &dto.ProductResponse{
		ID:     product.ID,
		Bundle: bundle,
	}, nil
}

func (s *storeService) GetReviews(ctx context.Context, q dto.ProductQuery) ([]json.RawMessage, error) {
	product, err := s.repo.GetProductCache(ctx, q.Shop, q.ProductID)
	if err != nil {
		return nil, err
	}
	return emptyIfNil(product.Reviews), nil
}

func normalizeStoreType(storeType string) string {
	if storeType == "" {
		return dto.DefaultStoreType
	}
	return storeType
}

func emptyIfNil(arr []json.RawMessage) []json.RawMessage {
	if arr == nil {
		return []json.RawMessage{}
	}
	return arr
}
