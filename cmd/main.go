package main

import (
	"fmt"
	"log"

	"github.com/jtaylorcpp/netquery/cmd/agent"
	"github.com/jtaylorcpp/netquery/cmd/query"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(agent.RunCmd)
	rootCmd.AddCommand(query.RunCmd)
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
