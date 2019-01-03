package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/nfgtech/broparser"
	"gitlab.com/nfgtech/netquery"
)

var broDir string
var postgresHost string
var postgresPort string
var postgresDB string
var postgresUser string
var postgresPass string

func init() {
	runCmd.Flags().StringVarP(&broDir, "bro-dir", "b", "", "bro log directory to watch")
	runCmd.Flags().StringVarP(&postgresHost, "psql-host", "", "localhost", "postgres hostname")
	runCmd.Flags().StringVarP(&postgresPort, "psql-port", "", "5432", "postgres port")
	runCmd.Flags().StringVarP(&postgresDB, "psql-db", "", "brodb", "postgres db name")
	runCmd.Flags().StringVarP(&postgresUser, "psql-user", "", "brouser", "postgres username")
	runCmd.Flags().StringVarP(&postgresPass, "psql-pass", "", "bropass", "postgres user pasword")

	rootCmd.AddCommand(runCmd)
}

func main() {
	execute()
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "netquery",
	Short: "netquery is a SQL service for network data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("welcome to netquery! Use the -h flag to see more options.")
	},
}

var runCmd = &cobra.Command{
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

			db := broparser.NewPostgresDB(postgresHost, postgresPort, postgresUser, postgresPass, postgresDB)

			netquery.StartAgent(tailers, db)
		}
	},
}
