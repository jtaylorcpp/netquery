package queries

import (
	"time"

	//"github.com/jinzhu/gorm"
	bp "github.com/jtaylorcpp/broparser"
)

/*
by ip address (client or server or both)
by port (client  or server or both)
by protocol
by service
by duration

before timestamp
after timestamp
between timesamp range

duration shorter than
duration longer than
duration in range

*/

// Returns all conn records in agiven database
func (q *query) BroConnAll() []bp.ConnRecord {
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

//
func (q *query) BroConn(clientAddrs []string, clientPorts []uint16,
	serverAddrs []string, serverPorts []uint16,
	protocols []string, services []string,
	beforeTs time.Time, afterTs time.Time) []bp.ConnRecord {

	tx := q.db

	for idx, caddrs := range clientAddrs {
		switch idx {
		case 0:
			tx = tx.Where("client_addr = ?", caddrs)
		default:
			tx = tx.Or("client_addr = ?", caddrs)
		}
	}

	for idx, cport := range clientPorts {
		switch idx {
		case 0:
			tx = tx.Where("client_port = ?", cport)
		default:
			tx = tx.Or("client_port = ?", cport)
		}
	}

	for idx, saddr := range serverAddrs {
		switch idx {
		case 0:
			tx = tx.Where("server_addr = ?", saddr)
		default:
			tx = tx.Or("server_addr = ?", saddr)
		}
	}

	for idx, sport := range serverPorts {
		switch idx {
		case 0:
			tx = tx.Where("server_port = ?", sport)
		default:
			tx = tx.Or("server_port = ?", sport)
		}
	}

	for idx, proto := range protocols {
		switch idx {
		case 0:
			tx = tx.Where("protocol = ?", proto)
		default:
			tx = tx.Or("protocol = ?", proto)
		}
	}

	for idx, service := range services {
		switch idx {
		case 0:
			tx = tx.Where("service = ?", service)
		default:
			tx = tx.Or("service = ?", service)
		}
	}

	if !beforeTs.IsZero() && !afterTs.IsZero() {
		// range
		tx = tx.Where("time_stamp >= ? AND time_stamp <= ?", afterTs, beforeTs)

	} else if !beforeTs.IsZero() && afterTs.IsZero() {
		// records before ts
		tx = tx.Where("time_stamp <= ?", beforeTs)

	} else if beforeTs.IsZero() && !afterTs.IsZero() {
		// records after ts
		tx = tx.Where("time_stamp >= ?", afterTs)

	} else {
		// no-op

	}

	return []bp.ConnRecord{}
}
