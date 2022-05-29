/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// migrateFromNotionCmd represents the migrateFromNotion command
var migrateFromNotionCmd = &cobra.Command{
	Use:   "migrateFromNotion",
	Short: "Pulls all movies out of Notion, builds local DB, and creates Grist document.",
	Run:   migrateFromNotion,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expected 2 arguments, got %v", len(args))
		}
		return nil
	},
}

func migrateFromNotion(cmd *cobra.Command, args []string) {
	// csvName := args[0]
	dbName := args[1]

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing db: %v", err)
		}
	}()
}

func init() {
	rootCmd.AddCommand(migrateFromNotionCmd)

}
