package gracenote

import (
	"testing"
	bolt"github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
)

func TestGestSports(t *testing.T) {

	os.Setenv("NEO4JBOLT", "bolt://neo4j:prelude@localhost:7687")
	GetSport("all")
	driver := bolt.NewDriver()
	conn, _ := driver.OpenNeo(os.Getenv("NEO4JBOLT"))

	result, _ := conn.ExecNeo("MATCH (n:Sport) return COUNT(n)", nil)

	rows, _ := result.RowsAffected()

	if rows == 0 {
		t.Error("Expected Results to be in the database but 0 were found")
	}

}