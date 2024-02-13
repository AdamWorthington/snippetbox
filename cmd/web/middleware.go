package main

import (
   "net/http"
   "fmt"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()
            if err != nil {
                w.Header().Set("Connection", "Close")

                app.serveError(w, r, fmt.Errorf("panic: %v", err))
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func (app *application) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        proto := r.Proto
        method := r.Method
        uri := r.URL.RequestURI()
        app.logger.Info("recieved request", "ip", ip, "proto", proto, "method", method, "uri", uri)
        next.ServeHTTP(w,r)
    })
}

func secureHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        w.Header().Set("Content-Security-Policy",
        "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
        w.Header().Set("Referrer-Policy", "origin-when-cross-origin") 
        w.Header().Set("X-Content-Type-Options", "nosniff") 
        w.Header().Set("X-Frame-Options", "deny") 
        w.Header().Set("X-XSS-Protection", "0")
        next.ServeHTTP(w, r)
    })
}