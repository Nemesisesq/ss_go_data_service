package dao

import (
	"testing"
	"fmt"
	"encoding/json"
	"os"
)


func setUp() {
	os.Setenv("NEO4JBOLT", "bolt://nem:prelude@0.0.0.0:7687")
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	//myTeardownFunction()
	os.Exit(retCode)
}

func TestGetSport(t *testing.T) {
	res := GetSport("70", nil)

	rJSON, _ := json.Marshal(res)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}