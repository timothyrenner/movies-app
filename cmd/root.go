/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var GRIST_KEY string

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

	gristKey, exists := os.LookupEnv("GRIST_KEY")
	if !exists {
		log.Println("Could not find GRIST_KEY in environment or .env.")
	}
	GRIST_KEY = gristKey
}
