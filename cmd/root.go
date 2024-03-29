/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var GRIST_KEY string
var GRIST_DOCUMENT_ID string
var OMDB_KEY string
var DB string = "./data/movies.db"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "movies-app",
	Short: "Creates and updates a movies database, and synchronizes with Grist.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.Println("Loading .env")
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file.")
	}

	omdbKey, exists := os.LookupEnv("OMDB_KEY")
	if !exists {
		log.Println("Could not find OMDB_KEY in environment or .env.")
	}
	OMDB_KEY = omdbKey
}
