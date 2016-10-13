package middleware

import (
	"gopkg.in/redis.v4"
	"github.com/gorilla/context"
	"net/http"
	"github.com/codegangsta/negroni"
)

type CacheAccessor struct {
	redis.Client
	url  string
	name string
	coll string
}

func NewCacheAccessor(url, name, coll string) (*CacheAccessor, error) {
	client, err := redis.NewClient()
	if err == nil {
		return &CacheAccessor{client , url, name, coll}, nil
	} else {
		return &CacheAccessor{}, err
	}
}

func (da *CacheAccessor) Set(request *http.Request, client redis.Client) {
	db := client.DB(da.name)
	context.Set(request, "db", db)
	context.Set(request, "mgoSession", client)
}

type Database struct {
	dba CacheAccessor
}

func NewDatabase(CacheAccessor CacheAccessor) *Database {
	return &Database{CacheAccessor}
}

func (d *Database) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		//reqSession := d.dba.Clone()
		//defer reqSession.Close()
		//d.dba.Set(request, reqSession)
		next(writer, request)
	}
}


