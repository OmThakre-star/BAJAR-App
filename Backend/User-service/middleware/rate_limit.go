package middleware

import (
"net/http"
"sync"
"time"


"golang.org/x/time/rate"
)


type client struct {
limiter *rate.Limiter
lastSeen time.Time
}


var clients = make(map[string]*client)
var mu sync.Mutex


func getLimiter(ip string) *rate.Limiter {
mu.Lock()
defer mu.Unlock()
c, exists := clients[ip]
if !exists {
lim := rate.NewLimiter(1, 5)
clients[ip] = &client{limiter: lim, lastSeen: time.Now()}
return lim
}
c.lastSeen = time.Now()
return c.limiter
}


func RateLimitMiddleware(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
ip := r.RemoteAddr
lim := getLimiter(ip)
if !lim.Allow() {
http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
return
}
next.ServeHTTP(w, r)
})
}