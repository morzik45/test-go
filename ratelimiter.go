package exam

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	ipCount map[string]int
	sync.Mutex
}

var limiter Limiter

func init() {
	limiter.ipCount = make(map[string]int)
}

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the IP address for the current user.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get the # of times the visitor has visited in the last 4 seconds
		limiter.Lock()
		count, ok := limiter.ipCount[ip]
		if !ok {
			limiter.ipCount[ip] = 0
		}
		if count > 20 {
			limiter.Unlock()
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		} else {
			limiter.ipCount[ip]++
		}
		time.AfterFunc(time.Second*4, func() {
			limiter.Lock()
			limiter.ipCount[ip]--
			limiter.Unlock()
		})
		if limiter.ipCount[ip] == 20 {
			// set it to 40 so the decrement timers will only decrease it to
			// 10, and they stay blocked until the next timer resets it to 0
			limiter.ipCount[ip] = 40
			time.AfterFunc(time.Minute*10, func() {
				limiter.Lock()
				limiter.ipCount[ip] = 0
				limiter.Unlock()
			})
		}
		limiter.Unlock()
		next.ServeHTTP(w, r)
	})
}
