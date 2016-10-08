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

// CreateTable ...
func CreateTable(tableName, token, serviceName string) string {
	query := `
  mutation {
    createTable(name: "{{.TableName}}") {
      db
    }
  }
  `
	s := struct {
		TableName string
	}{
		TableName: tableName,
	}
	queryString := tp.TemplateParse(query, s)

	graphql := grafetch.New("http://localhost:12345/graphql")
	graphql.SetHeader("x-meepcloud-access-token", token)
	graphql.SetHeader("x-meepcloud-service-name", serviceName)
	graphql.SetQuery(grafetch.GraphQLQuery{
		Query: queryString,
	})
	var resp struct {
		Data struct {
			CreateTable struct {
				DB string `json:"db"`
			} `json:"createTable"`
		} `json:"data"`
		Errors interface{} `json:"errors"`
	}
	graphql.Fetch(&resp)
	return resp.Data.CreateTable.DB
}

// CreateService ...
func CreateService(serviceName, accessToken string) string {
	query := `
  mutation {
    createService(name: "{{.ServiceName}}") {
      id
      serviceName
    }
  }
  `
	s := struct {
		ServiceName string
	}{
		ServiceName: serviceName,
	}
	queryString := tp.TemplateParse(query, s)

	graphql := grafetch.New("http://localhost:12345/graphql")
	graphql.SetHeader("x-meepcloud-access-token", accessToken)
	graphql.SetQuery(grafetch.GraphQLQuery{
		Query: queryString,
	})
	var resp struct {
		Data struct {
			CreateService struct {
				ID          string `json:"id"`
				ServiceName string `json:"serviceName"`
			} `json:"createService"`
		} `json:"data"`
		Errors interface{} `json:"errors"`
	}
	graphql.Fetch(&resp)
	return resp.Data.CreateService.ServiceName
}

// CreateAccount ...
func CreateAccount(username, password string) (accessToken string) {
	query := `
  mutation {
    createAccount(username: "{{.Username}}", password: "{{.Password}}") {
      accessToken
    }
  }
  `
	s := struct {
		Username string
		Password string
	}{
		Username: username,
		Password: password,
	}
	queryString := tp.TemplateParse(query, s)

	graphql := grafetch.New("http://localhost:12345/graphql")
	graphql.SetQuery(grafetch.GraphQLQuery{
		Query: queryString,
	})
	var resp struct {
		Data struct {
			CreateAccount struct {
				AccessToken string `json:"accessToken"`
			} `json:"createAccount"`
		} `json:"data"`
		Errors interface{} `json:"errors"`
	}
	graphql.Fetch(&resp)
	return resp.Data.CreateAccount.AccessToken
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
