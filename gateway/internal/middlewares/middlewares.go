package middlewares

import (
	context2 "context"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	"github.com/ce-final-project/backend_game_server/gateway/internal/service"
	httpErrors "github.com/ce-final-project/backend_game_server/pkg/http_errors"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/labstack/echo/v4"
	"strings"
	"time"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	AuthorizationMiddleware() echo.MiddlewareFunc
	//AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log logger.Logger
	cfg *config.Config
	as  service.AuthService
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config, as service.AuthService) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg, as: as}
}

func (mw *middlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start)

		if !mw.checkIgnoredURI(ctx.Request().RequestURI, mw.cfg.HTTP.IgnoreLogUrls) {
			mw.log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s)
		}

		return err
	}
}

//func (mw *middlewareManager) AuthorizationMiddleware() echo.MiddlewareFunc {
//
//	jwtConfig := middleware.JWTConfig{
//		TokenLookup: "header:" + mw.cfg.JWT.HeaderAuthorization,
//		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
//			keyFunc := func(t *jwt.Token) (interface{}, error) {
//				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
//					return nil, errors.New("jwt.Parse.Token.Method=" + t.Header["alg"].(string))
//				}
//				return []byte(mw.cfg.JWT.Secret), nil
//			}
//
//			// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
//			token, err := jwt.Parse(auth, keyFunc)
//			if err != nil {
//				return nil, err
//			}
//			if !token.Valid {
//				return nil, errors.New("invalid token")
//			}
//			return token, nil
//		},
//	}
//
//	return middleware.JWTWithConfig(jwtConfig)
//}

func (mw *middlewareManager) AuthorizationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context2.WithCancel(context2.Background())
			defer cancel()
			header := c.Request().Header
			token := header.Get(mw.cfg.JWT.HeaderAuthorization)
			mw.log.Debugf("Token from Header %v", token)
			result, err := mw.as.VerifyToken(ctx, token)
			if err != nil {
				return err
			}
			mw.log.Debugf("Result verifyToken %v", result)
			if !result.Valid {
				mw.log.WarnMsg("Invalid Token or Expired", httpErrors.WrongCredentials)
				return httpErrors.ErrorCtxResponse(c, httpErrors.WrongCredentials, mw.cfg.HTTP.DebugErrorsResponse)
			}
			c.Request().Header.Set("User-ID", result.AccountID.String())
			c.Request().Header.Set("Player-ID", result.PlayerID)
			return next(c)
		}
	}
}

func (mw *middlewareManager) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, s := range uriList {
		if strings.Contains(requestURI, s) {
			return true
		}
	}
	return false
}
