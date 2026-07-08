package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/dto"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/request"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/response"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/service"
)

type StoreController struct {
	store service.IStoreService
}

func NewStoreController(store service.IStoreService) *StoreController {
	return &StoreController{store: store}
}

func (c *StoreController) RegisterRoutes(r chi.Router) {
	r.Route("/storefront/stores", func(r chi.Router) {
		r.Get("/insight", c.GetInsights)
		r.Get("/reviews", c.GetReviews)
		r.Get("/pixels", c.GetPixels)
		r.Get("/product", c.GetProduct)
	})
}

func (c *StoreController) GetInsights(w http.ResponseWriter, r *http.Request) {
	var q dto.StoreQuery
	if err := request.Query(r, &q); err != nil {
		response.Error(w, err)
		return
	}

	data, err := c.store.GetInsights(r.Context(), q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, data)
}

func (c *StoreController) GetReviews(w http.ResponseWriter, r *http.Request) {
	var q dto.ProductQuery
	if err := request.Query(r, &q); err != nil {
		response.Error(w, err)
		return
	}

	data, err := c.store.GetReviews(r.Context(), q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, data)
}

func (c *StoreController) GetPixels(w http.ResponseWriter, r *http.Request) {
	var q dto.StoreQuery
	if err := request.Query(r, &q); err != nil {
		response.Error(w, err)
		return
	}

	data, err := c.store.GetPixels(r.Context(), q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, data)
}

func (c *StoreController) GetProduct(w http.ResponseWriter, r *http.Request) {
	var q dto.ProductQuery
	if err := request.Query(r, &q); err != nil {
		response.Error(w, err)
		return
	}

	data, err := c.store.GetProduct(r.Context(), q)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, data)
}
