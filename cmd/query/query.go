package query

import (
	"log"

	"github.com/jtaylorcpp/gerl/genserver"
	"github.com/jtaylorcpp/netquery/queries"
	"github.com/spf13/cobra"
)

var agentAddr string

func init() {
	RunCmd.Flags().StringVarP(&agentAddr, "netquery-agent", "a", "localhost:9001", "netquery agent address")

}

var RunCmd = &cobra.Command{
	Use:   "query",
	Short: "run a query against an agent node",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting the netquery query client.")
		if len(args) < 3 {
			log.Fatalf("3 or more args needed; args supplied: %v\n", args)
		}
		msg, err := queries.BuildQueryMSG(args[0], args[1], args[2:]...)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("message being sent to query server: %v\n", msg)
		returnMsg := genserver.Call(genserver.PidAddr(agentAddr), "query client", msg)
		log.Printf("recieved message back: %v\n", returnMsg)
	},
}
