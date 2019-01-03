package netquery

import (
	"log"

	"github.com/jtaylorcpp/broparser"
	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
)

func defaultCall(_ core.Pid, _ core.Message, _ genserver.FromAddr, s genserver.State) (core.Message, genserver.State) {
	return core.Message{}, s
}

func logCast(pid core.Pid, msg core.Message, faddr genserver.FromAddr, state genserver.State) genserver.State {
	log.Printf("genserver: %v\n", pid)
	log.Printf("cast msg: %v\n", msg)
	log.Printf("From addr: %v\n", faddr)
	log.Printf("current state: %v", state)

	s, ok := state.(ParserState)
	if !ok {
		log.Fatalf("unable to parse state to server staet: %v\n", state)
	}

	switch s.Type {
	case "conn.log":
		parsed := broparser.ParseConn(msg.GetValues()[0])
		if parsed.UID != "" {
			log.Printf("conn record parsed: %v\n", parsed)
			s.DB.Create(&parsed)
		}
	case "dns.log":
		parsed := broparser.ParseDNS(msg.GetValues()[0])
		if parsed.UID != "" {
			log.Printf("dns record parsed: %v\n", parsed)
			s.DB.Create(&parsed)
		}
	}

	return state
}
