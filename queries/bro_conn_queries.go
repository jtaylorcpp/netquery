package queries

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	//"github.com/jinzhu/gorm"
	bp "github.com/jtaylorcpp/broparser"
)

type BroConnQuery struct {
	ClientAddrs []string  `json:"client_addrs"`
	ClientPorts []uint16  `json:"client_ports"`
	ServerAddrs []string  `json:"server_addrs"`
	ServerPorts []uint16  `json:"server_ports"`
	Protocols   []string  `json:"Protocols"`
	Services    []string  `json:"Services"`
	BeforeTs    time.Time `json:"before_ts"`
	AfterTs     time.Time `json:"after_ts"`
}

func BroConnQueryFromKV(fields map[string]string) QueryInterface {

	query := BroConnQuery{
		ClientAddrs: []string{},
		ClientPorts: []uint16{},
		ServerAddrs: []string{},
		ServerPorts: []uint16{},
		Protocols:   []string{},
		Services:    []string{},
		BeforeTs:    time.Time{},
		AfterTs:     time.Time{},
	}

	for k, v := range fields {
		switch k {
		case "client_addrs":
			query.ClientAddrs = strings.Split(v, ",")
		case "client_ports":
		case "server_addrs":
			query.ServerAddrs = strings.Split(v, ",")
		case "server_ports":
		case "Protocols":
			query.Protocols = strings.Split(v, ",")
		case "Services":
			query.Services = strings.Split(v, ",")
		case "ts_after":
		case "ts_before":
		default:
			log.Printf("unknown conn key %v with value %v\n", k, v)
		}

	}
	return query
}

func BroConnQueryFromJSON(jsonBytes []byte) QueryInterface {
	broConn := BroConnQuery{}
	json.Unmarshal(jsonBytes, &broConn)

	return broConn
}

/*
Queries
*/

// Returns all conn records in agiven database
func (q *Query) BroConnAll() []bp.ConnRecord {
	rows, err := q.db.Exec("SELECT * FROM conn_records;").Rows()
	if err != nil {
		panic(err)
	}

	records := []bp.ConnRecord{}

	for rows.Next() {
		var record bp.ConnRecord
		q.db.ScanRows(rows, record)
		records = append(records, record)
	}

	return records
}

//Builds Query for BroConn based with flexible params
func (q *Query) BroConn(ClientAddrs []string, ClientPorts []uint16,
	ServerAddrs []string, ServerPorts []uint16,
	Protocols []string, Services []string,
	BeforeTs time.Time, AfterTs time.Time) []bp.ConnRecord {

	tx := q.db

	for idx, caddrs := range ClientAddrs {
		switch idx {
		case 0:
			tx = tx.Where("client_addr = ?", caddrs)
		default:
			tx = tx.Or("client_addr = ?", caddrs)
		}
	}

	for idx, cport := range ClientPorts {
		switch idx {
		case 0:
			tx = tx.Where("client_port = ?", cport)
		default:
			tx = tx.Or("client_port = ?", cport)
		}
	}

	for idx, saddr := range ServerAddrs {
		switch idx {
		case 0:
			tx = tx.Where("server_addr = ?", saddr)
		default:
			tx = tx.Or("server_addr = ?", saddr)
		}
	}

	for idx, sport := range ServerPorts {
		switch idx {
		case 0:
			tx = tx.Where("server_port = ?", sport)
		default:
			tx = tx.Or("server_port = ?", sport)
		}
	}

	for idx, proto := range Protocols {
		switch idx {
		case 0:
			tx = tx.Where("protocol = ?", proto)
		default:
			tx = tx.Or("protocol = ?", proto)
		}
	}

	for idx, service := range Services {
		switch idx {
		case 0:
			tx = tx.Where("service = ?", service)
		default:
			tx = tx.Or("service = ?", service)
		}
	}

	if !BeforeTs.IsZero() && !AfterTs.IsZero() {
		// range
		tx = tx.Where("time_stamp >= ? AND time_stamp <= ?", AfterTs, BeforeTs)

	} else if !BeforeTs.IsZero() && AfterTs.IsZero() {
		// records before ts
		tx = tx.Where("time_stamp <= ?", BeforeTs)

	} else if BeforeTs.IsZero() && !AfterTs.IsZero() {
		// records after ts
		tx = tx.Where("time_stamp >= ?", AfterTs)

	} else {
		// no-op

	}

	return []bp.ConnRecord{}
}
