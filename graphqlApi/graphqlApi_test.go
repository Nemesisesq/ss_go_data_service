package graphqlApi

import (
	"testing"
	"github.com/graphql-go/graphql"
	"fmt"
	log"github.com/Sirupsen/logrus"
	"encoding/json"
)

func TestGraph(t *testing.T)  {

	schema := Schema()

	query := `
        {
            hello
        }
    `
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphqllib operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}
func TestExecuteQuery() {

}