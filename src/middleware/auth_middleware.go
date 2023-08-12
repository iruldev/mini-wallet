package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iruldev/mini-wallet/src/constant"
	"github.com/iruldev/mini-wallet/src/helper"
	"github.com/iruldev/mini-wallet/src/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "token"
)

func AuthMiddleware(tokenMaker token.Maker) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := helper.PlugResponse(w)

			authorizationHeader := r.Header.Get(authorizationHeaderKey)

			if len(authorizationHeader) == 0 {
				err := errors.New("authorization header is not provided")
				_ = res.ReplyCustom(http.StatusUnauthorized, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				_ = res.ReplyCustom(http.StatusUnauthorized, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
				return
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationTypeBearer {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				_ = res.ReplyCustom(http.StatusUnauthorized, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
				return
			}

			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				_ = res.ReplyCustom(http.StatusUnauthorized, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), constant.AuthorizationPayloadKey, payload))
			next.ServeHTTP(w, r)
		})
	}
}
