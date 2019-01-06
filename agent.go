package netquery

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hpcloud/tail"
	"github.com/jinzhu/gorm"
	"github.com/jtaylorcpp/gerl/core"
	"github.com/jtaylorcpp/gerl/genserver"
	"github.com/jtaylorcpp/gerl/registrar"
	"github.com/jtaylorcpp/netquery/queries"
)

func StartAgent(logfiles map[string]*tail.Tail, db *gorm.DB, registrarPort string, queryPort string) {

	reg := registrar.NewRegistrarWithPort(core.GlobalScope, registrarPort)
	log.Printf("Registrar started at: %v", reg.Pid.GetAddr())

	servers := []*genserver.GenServer{}

	log.Println("Starting query service")
	querygs := genserver.NewGenServerWithPort(queries.NewQuery(db), core.GlobalScope, queryPort, postgresCall, postgresCast)
	servers = append(servers, querygs)
	queryRecord := registrar.NewRecord("query", querygs.Pid.GetAddr(), core.GlobalScope)
	log.Printf("dding record to registrar: %v\n", queryRecord)
	if registrar.AddRecords(reg.Pid.GetAddr(), "query server", queryRecord) {
		log.Printf("genserver %s added to registrar\n", "query")
	} else {
		log.Fatalf("genserver %s failed to be added to registrar\n", "query")
	}

	for logname, tailer := range logfiles {
		loggs := genserver.NewGenServer(ParserState{Type: logname, DB: db.New()}, core.LocalScope, defaultCall, logCast)

		log.Printf("genserver for %s at addr %v\n", logname, loggs.Pid.GetAddr())

		logRec := registrar.NewRecord(logname, loggs.Pid.GetAddr(), core.LocalScope)
		log.Printf("adding record to registrar: %v\n", logRec)

		if registrar.AddRecords(reg.Pid.GetAddr(), "filewatcher server init", logRec) {
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
