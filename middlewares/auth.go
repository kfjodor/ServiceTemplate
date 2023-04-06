package middlewares

import (
	"context"
	"net/http"

	"github.com/jinzhu/gorm"

	"service_template/logger"
)

func AuthMiddlewareGenerator(ctx context.Context, db *gorm.DB) (mw func(http.Handler) http.Handler) {
	log := logger.FromContext(ctx).WithField("m", "AuthMiddlewareGenerator")
	log.Debugf("AuthMiddlewareGenerator:: ")

	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/*
				tokenHeader := r.Header.Get("Authorization") // get token from HTTP header

				if tokenHeader == "" {
					handlers.ERROR_AUTH_MISSING(w)

					return
				}

				adminUser := new(models.AdminUser)
				db.First(&adminUser, "password = ?", tokenHeader)
				if adminUser.ID == 0 {
					handlers.ERROR_AUTH_INVALID(w, tokenHeader)

					return
				}
			*/
			next.ServeHTTP(w, r)
		})
	}

	return
}
