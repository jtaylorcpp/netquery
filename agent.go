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
)

func StartAgent(logfiles map[string]*tail.Tail, db *gorm.DB) {

	reg := registrar.NewRegistrar(core.LocalScope)
	log.Printf("Registrar started at: %v", reg.Pid.GetAddr())

	servers := []*genserver.GenServer{}

	for logname, tailer := range logfiles {
		loggs := genserver.NewGenServer(ParserState{Type: logname, DB: db.New()}, core.LocalScope, defaultCall, logCast)

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
