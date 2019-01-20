package agent

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hpcloud/tail"
	"github.com/jtaylorcpp/netquery"
	"github.com/jtaylorcpp/netquery/databases"
	"github.com/spf13/cobra"
)

var broDir string
var postgresHost string
var postgresPort string
var postgresDB string
var postgresUser string
var postgresPass string
var queryAddr string
var gerlPort string

func init() {
	RunCmd.Flags().StringVarP(&broDir, "bro-dir", "b", "", "bro log directory to watch")
	RunCmd.Flags().StringVarP(&postgresHost, "psql-host", "", "localhost", "postgres hostname")
	RunCmd.Flags().StringVarP(&postgresPort, "psql-port", "", "5432", "postgres port")
	RunCmd.Flags().StringVarP(&postgresDB, "psql-db", "", "brodb", "postgres db name")
	RunCmd.Flags().StringVarP(&postgresUser, "psql-user", "", "brouser", "postgres username")
	RunCmd.Flags().StringVarP(&postgresPass, "psql-pass", "", "bropass", "postgres user pasword")
	RunCmd.Flags().StringVarP(&queryAddr, "netquery-port", "p", "9001", "netquery port")
	RunCmd.Flags().StringVarP(&gerlPort, "gerl-registrar-port", "r", "9000", "gerl registrar port")

	agentErase.Flags().StringVarP(&postgresHost, "psql-host", "", "localhost", "postgres hostname")
	agentErase.Flags().StringVarP(&postgresPort, "psql-port", "", "5432", "postgres port")
	agentErase.Flags().StringVarP(&postgresDB, "psql-db", "", "brodb", "postgres db name")
	agentErase.Flags().StringVarP(&postgresUser, "psql-user", "", "brouser", "postgres username")
	agentErase.Flags().StringVarP(&postgresPass, "psql-pass", "", "bropass", "postgres user pasword")

	RunCmd.AddCommand(agentErase)
}

var RunCmd = &cobra.Command{
	Use:   "agent",
	Short: "run netquery in agent mode",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting the netquery agent.")
		if broDir != "" {
			broInfo, err := os.Stat(broDir)
			if err != nil {
				panic(err)
			}

			var path = ""
			var files = []string{}

			if !broInfo.IsDir() {
				//log.Printf("Watching file: %v\n", broDir)
				abspath, err := filepath.Abs(broDir)
				if err != nil {
					panic(err)
				}
				path = filepath.Dir(abspath)
				files = append(files, filepath.Base(broDir))
			} else {
				//log.Printf("Watching directory: %v\n", broDir)
				abspath, err := filepath.Abs(broDir)
				if err != nil {
					panic(err)
				}
				path = abspath
				filesDir, err := ioutil.ReadDir(path)
				if err != nil {
					panic(err)
				}

				for _, f := range filesDir {
					files = append(files, f.Name())
				}
			}

			//log.Printf("path: %v\n", path)
			//log.Printf("files: %v\n", files)
			tailers := netquery.GetBroTailers(path, files)
			log.Println(tailers)

			brodb := databases.NewPostgresDB(postgresHost, postgresPort, postgresUser, postgresPass, postgresDB)
			databases.PSQLInit(brodb)
			defer brodb.Close()

			netquery.StartAgent(tailers, brodb, gerlPort, queryAddr)
		} else {
			brodb := databases.NewPostgresDB(postgresHost, postgresPort, postgresUser, postgresPass, postgresDB)
			databases.PSQLInit(brodb)
			defer brodb.Close()

			netquery.StartAgent(map[string]*tail.Tail{}, brodb, gerlPort, queryAddr)
		}
	},
}

var agentErase = &cobra.Command{
	Use:   "erase",
	Short: "erase all data collected by agent",
	Run: func(cmd *cobra.Command, args []string) {
		brodb := databases.NewPostgresDB(postgresHost, postgresPort, postgresUser, postgresPass, postgresDB)
		databases.PSQLClear(brodb)
		brodb.Close()
	},
}
