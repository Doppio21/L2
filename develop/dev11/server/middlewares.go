package server

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(log *logrus.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			next.ServeHTTP(w, r)
			log.Infof("| %s | %s | %s |",
				time.Since(now),
				r.Method,
				r.URL.String(),
			)
		}
		return http.HandlerFunc(fn)
	}
}
