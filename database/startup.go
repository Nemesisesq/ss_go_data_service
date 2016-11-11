package database

import (
	"github.com/joeljames/nigroni-mgo-session"
	com "github.com/nemesisesq/ss_data_service/common"
	"os"
)

func DBStartup() nigronimgosession.DatabaseAccessor {
	dbURL := os.Getenv("MONGODB_URI")
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
