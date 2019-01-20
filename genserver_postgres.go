package netquery

import (
	"log"

	//"github.com/jtaylorcpp/broparser"
	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
)

func postgresCall(_ core.Pid, msg core.Message, faddr genserver.FromAddr, s genserver.State) (core.Message, genserver.State) {
	log.Printf("postgres call from %v with message %v\n", faddr, msg)

	switch msg.Type {
	case core.Message_SIMPLE:
		switch msg.Subtype {
		case core.Message_GET:
			//parse logs
			log.Println("postgres query: ", msg)
			return msg, s
		default:
			log.Printf("postgres does not handle subtype %v\n", msg.Subtype)
		}

	default:
		log.Printf("postgres does not handle msg type %v\n", msg.Type)
	}
	return core.Message{}, s
}

func postgresCast(_ core.Pid, _ core.Message, _ genserver.FromAddr, s genserver.State) genserver.State {
	return s
}
