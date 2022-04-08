package middleware

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

type IMiddleware interface {
	PanicRecovery(next http.Handler) http.Handler
	RequestLogger(next http.Handler) http.Handler
}

type middleware struct {
	logger *zap.Logger
}

func NewMiddleware(logger *zap.Logger) IMiddleware {
	middleware := new(middleware)
	middleware.logger = logger
	return middleware
}

func (m *middleware) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				m.logger.Error("PanicRecovery", zap.String("Stack", string(debug.Stack())))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func (m *middleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		m.logger.Info("RequestLogger", zap.String("Method", req.Method), zap.String("RequestURI", req.RequestURI), zap.String("RemoteAddr", req.RemoteAddr))
		next.ServeHTTP(w, req)
	})
}
