package netquery

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hpcloud/tail"
	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
	"github.com/jtaylorcpp/gerl/registrar"
	"gitlab.com/nfgtech/broparser"
)

type ParserState struct {
	Type string
}

func StartServer(logfiles map[string]*tail.Tail) {

	reg := registrar.NewRegistrar(core.LocalScope)
	log.Printf("Registrar started at: %v", reg.Pid.GetAddr())

	servers := []*genserver.GenServer{}

	for logname, tailer := range logfiles {
		loggs := genserver.NewGenServer(ParserState{Type: logname}, core.LocalScope, defaultCall, logCast)

		log.Printf("genserver for %s at addr %v\n", logname, loggs.Pid.GetAddr())

		logRec := registrar.NewRecord(logname, loggs.Pid.GetAddr(), core.LocalScope)
		log.Printf("adding record to registrar: %v\n", logRec)

		if registrar.AddRecords(reg.Pid.GetAddr(), "server init", logRec) {
			log.Printf("genserver %s added to registrar\n", logname)
		} else {
			log.Fatalf("genserver %s cannot be added to registrar", logname)
		}

		go func() {
			defer tailer.Cleanup()
			for line := range tailer.Lines {
				lineMsg := core.Message{
					Type:        core.Message_SIMPLE,
					Description: logname,
					Values:      []string{line.Text},
				}

				localaddr := fmt.Sprintf("tailer:%s", logname)
				genserver.Cast(genserver.PidAddr(loggs.Pid.GetAddr()), genserver.PidAddr(localaddr), lineMsg)
			}
		}()

		servers = append(servers, loggs)
	}

	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Interrupt)
	signal.Notify(killChan, os.Kill)

	<-killChan

	reg.Terminate()

	for _, server := range servers {
		server.Terminate()
		log.Printf("genserver %v terminated with state %v\n", server.Pid.GetAddr(), server.State)
	}

}

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
		}
	case "dns.log":
		parsed := broparser.ParseDNS(msg.GetValues()[0])
		if parsed.UID != "" {
			log.Printf("dns record parsed: %v\n", parsed)
		}
	}

	return state
}
