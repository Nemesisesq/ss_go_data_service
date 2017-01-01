package middleware

import (
	"context"
	"net/http"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/redis.v5"
)

type CacheAccessor struct {
	client redis.Client
	addr   string
	pass   string
	db     int
}

func NewCacheAccessor(addr, pass string, db int) (*CacheAccessor, error) {
	logrus.WithFields(logrus.Fields{
		"address":  addr,
		"password": pass,
		"db":       db,
	}).Info("This is the information to connect to redis")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	logrus.Info(addr)

	pong, err := client.Ping().Result()
	common.Check(err)
	fmt.Printf("redis %v", pong)

	return &CacheAccessor{*client, addr, pass, db}, nil
}

func (ca *CacheAccessor) Set(request *http.Request, client redis.Client) context.Context {
	//gcontext.Set(request, "client", *client)
	return context.WithValue(request.Context(), "redis_client", &client)
}

type RedisClient struct {
	rca CacheAccessor
}

func NewRedisClient(CacheAccessor CacheAccessor) *RedisClient {
	return &RedisClient{CacheAccessor}
}

func (ca *RedisClient) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		ctx := ca.rca.Set(request, ca.rca.client)
		next(writer, request.WithContext(ctx))
	}
}
