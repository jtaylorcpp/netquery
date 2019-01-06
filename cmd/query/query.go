package query

import (
	"log"

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
		log.Println("Starting the netquery client.")

		/*brodb := databases.NewPostgresDB(postgresHost, postgresPort, postgresUser, postgresPass, postgresDB)
		databases.PSQLInit(brodb)
		defer brodb.Close()

		netquery.StartAgent(tailers, brodb)*/
	},
}
