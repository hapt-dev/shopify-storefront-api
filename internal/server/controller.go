package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/app"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/controller"
)

// Controllers groups the HTTP controllers of every module.
type Controllers struct {
	Store *controller.StoreController
}

// newControllers builds the controllers from the services initialized in app.
func newControllers(services *app.Services) *Controllers {
	return &Controllers{
		Store: controller.NewStoreController(services.Store),
	}
}

// registerRoutes mounts each controller's routes onto the router.
func (c *Controllers) registerRoutes(r chi.Router) {
	c.Store.RegisterRoutes(r)
}
