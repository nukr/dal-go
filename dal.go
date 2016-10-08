package dal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nukr/grafetch"
	"github.com/nukr/template_parse"
)

// DAL ...
type DAL struct {
	URL         string
	AccessToken string
	ServiceName string
}

// CreateObject ...
func (d DAL) CreateObject(tableName string, variables interface{}) interface{} {
	query := `
  mutation Q($body: PlainObject){
    createObject(table: "{{.TableName}}", body: $body)
  }
  `
	s := struct {
		TableName string
	}{
		TableName: tableName,
	}

	queryString := tp.TemplateParse(query, s)
	g := grafetch.New(d.URL)
	g.SetHeader("x-meepcloud-access-token", d.AccessToken)
	g.SetHeader("x-meepcloud-service-name", d.ServiceName)
	bs, err := json.Marshal(struct {
		Body interface{} `json:"body"`
	}{
		Body: variables,
	})
	fmt.Println(string(bs))
	if err != nil {
		log.Fatal(err)
	}
	g.SetQuery(grafetch.GraphQLQuery{
		Query:     queryString,
		Variables: string(bs),
	})
	var result interface{}
	err = g.Fetch(&result)
	if err != nil {
		log.Fatal(err)
	}
	if err := g.Fetch(&result); err != nil {
		log.Fatal(err)
	}
	return result
}
