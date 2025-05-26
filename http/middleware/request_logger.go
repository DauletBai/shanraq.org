// github.com/DauletBai/shanraq.org/http/middleware/request_logger.go
package middleware 

import (
    "net/http"
    "time"

    "github.com/go-chi/chi/v5/middleware" 
    k "github.com/DauletBai/shanraq.org/core/kernel"
)

func RequestLogger(kernel *k.Kernel) func(next http.Handler) http.Handler {
    logger := kernel.Logger()

    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor) 
            tStart := time.Now()

            reqID := middleware.GetReqID(r.Context())

            defer func() {
                status := ww.Status()
                duration := time.Since(tStart)

                logFields := []interface{}{
                    "method", r.Method,
                    "path", r.URL.Path,
                    "status", status,
                    "duration_ms", float64(duration.Nanoseconds()) / float64(time.Millisecond),
                    "remote_addr", r.RemoteAddr,
                    "user_agent", r.UserAgent(),
                }
                if reqID != "" {
                    logFields = append(logFields, "request_id", reqID)
                }

                if status >= 400 && status < 500 { 
                    logger.Warn("Client error request", logFields...)
                } else if status >= 500 {
                    logger.Error("Server error request", logFields...)
                } else { 
                    logger.Info("Handled request", logFields...)
                }
            }()

            next.ServeHTTP(ww, r)
        }
        return http.HandlerFunc(fn)
    }
}