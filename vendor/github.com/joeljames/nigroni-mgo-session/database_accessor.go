package nigronimgosession

import (
	"net/http"

	//"github.com/gorilla/context"
	"context"
	mgo "gopkg.in/mgo.v2"
	"github.com/Sirupsen/logrus"
)

type DatabaseAccessor struct {
	*mgo.Session
	url  string
	name string
	coll string
}

func NewDatabaseAccessor(url, name, coll string) (*DatabaseAccessor, error) {


	session, err := mgo.Dial(url)
	logrus.Error(err)
	if err == nil {
		return &DatabaseAccessor{session, url, name, coll}, nil
	} else {
		return &DatabaseAccessor{}, err
	}
}

func (da *DatabaseAccessor) Set(request *http.Request, session *mgo.Session) context.Context {
	db := session.DB(da.name)

	return context.WithValue(request.Context(), "db", *db)
}
