package database

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joeljames/nigroni-mgo-session"
	com "github.com/nemesisesq/ss_data_service/common"
)

func DBStartup() nigronimgosession.DatabaseAccessor {

	var f string
	if os.Getenv("MONGODB_URI") != "" {
		f = os.Getenv("MONDODB_URI")
		logrus.Info("f", f)
	} else {

		f = os.Getenv("MONGODB_PORT_27017_TCP_ADDR")
	}

	dbURL := f
	// Use the MongoDB `DATABASE_NAME` from the env
	dbName := GetDatabase()
	// Set the MongoDB collection name
	dbColl := GetCollection()

	com.AnnounceMongoConnection(dbURL, dbName, dbColl)

	// Creating the database accessor here.
	// Pointer to this database accessor will be passed to the middleware.
	dbAccessor, err := nigronimgosession.NewDatabaseAccessor(dbURL, dbName, dbColl)

	com.Check(err)

	return *dbAccessor
}
