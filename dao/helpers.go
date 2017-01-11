package dao

import (
	"os"
	"github.com/nemesisesq/ss_data_service/common"
	"github.com/Sirupsen/logrus"
	bolt"github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func QueryNeo(cypher_query string, params map[string]interface{}) ([][]interface{}, map[string]interface{}) {
	driver := bolt.NewDriver()
	logrus.Info(os.Getenv("NEO4JBOLT"))
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	common.Check(err)
	defer conn.Close()

	data, rowMetaData, _, _ := conn.QueryNeoAll(cypher_query, params)

	return data, rowMetaData
}

func ProcessCypher(conn bolt.Conn, cypher_template string, params map[string]interface{}) {
	stmt, err := conn.PrepareNeo(cypher_template)
	common.Check(err)
	result, err := stmt.ExecNeo(params)
	common.Check(err)
	logrus.WithFields(logrus.Fields{
		"Cypher QueryResult": result,
		"params":             params,
	}).Info()
	if val, ok := result.RowsAffected(); ok != nil {
		//numResult, err := result.RowsAffected()
		common.Check(err)
		logrus.Info("CREATED ROWS", val)

	}
	//Closing the statment will also close the rows
	stmt.Close()
}
