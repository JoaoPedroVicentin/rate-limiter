package main

import (
	"net/http"
	"time"

	. "github.com/JoaoPedroVicentin/rate-limiter/configs"
	. "github.com/JoaoPedroVicentin/rate-limiter/internal/entity/client_entity"
	. "github.com/JoaoPedroVicentin/rate-limiter/internal/entity/rate_limiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

func main() {
	configs, err := LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     configs.RedisAddress,
		Password: configs.RedisPassword,
	})

	redisClient := NewRedisClient(client)

	rateLimiter := NewRateLimiter(
		redisClient,
		configs.MaxRequestsPerIp,
		configs.MaxRequestsPerApiKey,
		time.Duration(configs.TimeInterval)*time.Minute,
		time.Duration(configs.BlockDuration)*time.Minute,
	)

	router := chi.NewRouter()
	router.Handle("/", http.HandlerFunc(HelloWorld))

	handler := RateLimiterMiddleware(rateLimiter, router, configs.ApiKey)

	http.ListenAndServe(":"+configs.WebServerPort, handler)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
