package queries

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/jtaylorcpp/gerl/core"
)

// queries local db and returns json result
func QueryFromMsg(db *gorm.DB, msg core.Message) ([]byte, error) {
	query := NewQuery(db)

	if len(msg.Values) < 3 {
		return []byte{}, errors.New(fmt.Sprintf("not enough arguments to query: %v\n", msg.Values))
	}

	source, table, args := msg.Values[0], msg.Values[1], msg.Values[2:]

	if _, ok := jsonLogParser[source]; ok {
		if _, tableOk := jsonLogParser[source][table]; tableOk {

			switch args[0] {
			case "all":
				query := NewQuery(db)
				broConn := query.BroConnAll()
				log.Printf("queried data")
				return json.Marshal(broConn)
			default:
				query := NewQuery(db)
				broConnJSON := jsonLogParser[source][table]([]byte(args[0]))
				switch broConnJSON.(type) {
				case BroConnQuery:
					log.Printf("query to run: %v\n", string(broConnJSON))
					return []byte{}, nil
				}
			}
		} else {
			return core.Message{}, errors.New(fmt.Sprintf("table %v does not exist in source %v to query\n", table, source))
		}
	} else {
		return core.Message{}, errors.New(fmt.Sprintf("source %v does not exist to query\n", source))
	}
}
