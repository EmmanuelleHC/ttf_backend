package middlewares

import (
	"net/http"

	"github.com/Aguztinus/petty-cash-backend/constants"
	"github.com/Aguztinus/petty-cash-backend/models/dto"

	"github.com/Aguztinus/petty-cash-backend/api/services"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/pkg/echox"
	"github.com/labstack/echo/v4"
)

// CorsMiddleware middleware for cors
type CasbinMiddleware struct {
	handler lib.HttpHandler
	logger  lib.Logger
	config  lib.Config

	casbinService services.CasbinService
}

// NewCorsMiddleware creates new cors middleware
func NewCasbinMiddleware(
	handler lib.HttpHandler,
	logger lib.Logger,
	config lib.Config,
	casbinService services.CasbinService,
) CasbinMiddleware {
	return CasbinMiddleware{
		handler:       handler,
		logger:        logger,
		config:        config,
		casbinService: casbinService,
	}
}

func (a CasbinMiddleware) core() echo.MiddlewareFunc {
	prefixes := a.config.Casbin.IgnorePathPrefixes

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			request := ctx.Request()
			if isIgnorePath(request.URL.Path, prefixes...) {
				return next(ctx)
			}

			p := ctx.Request().URL.Path
			m := ctx.Request().Method
			claims, ok := ctx.Get(constants.CurrentUser).(*dto.JwtClaims)
			if !ok {
				return echox.Response{Code: http.StatusUnauthorized}.JSON(ctx)
			}

			if ok, err := a.casbinService.Enforcer.Enforce(claims.ID, p, m); err != nil {
				return echox.Response{Code: http.StatusForbidden, Message: err}.JSON(ctx)
			} else if !ok {
				return echox.Response{Code: http.StatusForbidden}.JSON(ctx)
			}

			return next(ctx)
		}
	}
}

func (a CasbinMiddleware) Setup() {
	if !a.config.Casbin.Enable {
		return
	}

	a.logger.Zap.Info("Setting up casbin middleware")
	a.handler.Engine.Use(a.core())
}
