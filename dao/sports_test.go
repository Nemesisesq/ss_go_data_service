package dao

import (
	"testing"
	"fmt"
	"encoding/json"
	"os"
)


func setUp() {
	os.Setenv("NEO4JBOLT", "bolt://nem:prelude@0.0.0.0:7687")
	bolt_url = os.Getenv("NEO4JBOLT")
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	//myTeardownFunction()
	os.Exit(retCode)
}

func TestGetSports(t *testing.T) {
	res := GetSport(nil)

	rJSON, _ := json.Marshal(res)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}

func TestGetTeams(t *testing.T) {

	res := GetTeams("111")

	rJSON, _  := json.Marshal(res)
	fmt.Printf("%s \n", rJSON)
}

func TestGetOrgs(t *testing.T) {
	res := GetOrgs("111")
	rJSON, _  := json.Marshal(res)
	fmt.Printf("%s \n", rJSON)
}