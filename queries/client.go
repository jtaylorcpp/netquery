package queries

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jtaylorcpp/gerl/core"
)

//
func BuildQueryMSG(source string, table string, args ...string) (core.Message, error) {
	if _, ok := kvLogParser[source]; ok {
		if _, ok := kvLogParser[source][table]; ok {
			msg := core.Message{
				Type:    core.Message_SIMPLE,
				Subtype: core.Message_GET,
				Values: []string{
					source,
					table,
				},
			}

			switch len(args) {
			case 0:
				return core.Message{}, errors.New("no query args passed")
			case 1:
				if args[0] == "all" {
					msg.Values = append(msg.Values, "all")
				} else {
					return core.Message{}, errors.New(fmt.Sprintf("single query arg given that is not 'all': %v\n", args))
				}
			default:
				// turn args into key value
				kv := map[string]string{}
				for i := 0; i < (len(args) - (len(args) % 2)); i += 2 {
					kv[args[i]] = args[i+1]
				}
				log.Printf("query kv to parse: %v\n", kv)
				parsedQuery := kvLogParser[source][table](kv)
				//var queryJSON string
				switch parsedQuery.(type) {
				case BroConnQuery:
					query, ok := parsedQuery.(BroConnQuery)
					if !ok {
						log.Fatal("cannont convert query to struct")
					}
					log.Printf("parsed query: %#v\n", query)
					queryJson, err := json.Marshal(query)
					log.Printf("parsed query json: %v\n", string(queryJson))
					if err != nil {
						return core.Message{}, err
					}
					msg.Values = append(msg.Values, string(queryJson))
				}

			}

			log.Printf("parsed msg: %v\n", msg)
			return msg, nil

		} else {
			return core.Message{}, errors.New(fmt.Sprintf("log source not supported; supported types are: %v\n", kvLogParser))
		}
	} else {
		return core.Message{}, errors.New(fmt.Sprintf("log source not supported; supported types are: %v\n", kvLogParser))
	}
}
