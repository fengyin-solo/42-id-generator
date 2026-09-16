// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"idgenerator/internal/config"
	"idgenerator/internal/model"
	"idgenerator/internal/service"
	"idgenerator/internal/store"
	"idgenerator/pkg/httpx"
	"idgenerator/pkg/logger"
)

type Server struct {
	svc   *service.Service
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

func NewServer(svc *service.Service, store store.Store, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, store: store, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerBizTypeRoutes(mux)
	s.registerMachineNodeRoutes(mux)
	s.registerIDRuleRoutes(mux)
	s.registerSegmentRoutes(mux)
	s.registerAllocRecordRoutes(mux)
	s.registerLeaseRoutes(mux)
	s.registerNodeHeartbeatRoutes(mux)
	s.registerSnowflakeConfigRoutes(mux)
	s.registerAllocStatsRoutes(mux)
	s.registerRecycleRecordRoutes(mux)
	s.registerIDGenRoutes(mux)
	s.registerExportImportRoutes(mux)
	// 静态文件
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.loggingMiddleware(s.recoveryMiddleware(s.authMiddleware(s.rateLimitMiddleware(mux))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 静态文件和根路径跳过鉴权
		if r.URL.Path == "/" || r.URL.Path == "/index.html" ||
			r.URL.Path == "/style.css" || r.URL.Path == "/app.js" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-Api-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if s.cfg != nil && key != s.cfg.APIKey {
			httpx.Unauthorized(w, "无效的 API Key")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	bucket := make(chan struct{}, 100)
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case bucket <- struct{}{}:
			default:
			}
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-bucket:
			next.ServeHTTP(w, r)
		default:
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
		}
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
