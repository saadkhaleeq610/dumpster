/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting PostgreSQL backup...")

		cmdArgs := []string{
			"-h", host,
			"-p", strconv.Itoa(port),
			"-U", user,
			"-F", "c",
			"-f", output,
			dbname,
		}

		dumpCmd := exec.Command("pg_dump", cmdArgs...)

		dumpCmd.Env = append(os.Environ(), "PGPASSWORD="+password)

		dumpCmd.Stdout = os.Stdout
		dumpCmd.Stderr = os.Stderr

		err := dumpCmd.Run()
		if err != nil {
			log.Fatalf("Backup failed: %v", err)
		}

		fmt.Printf("Backup complete: %s\n", output)
	},
}

var (
	host     string
	port     int
	user     string
	password string
	dbname   string
	output   string
)

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVar(&host, "host", "localhost", "Database host")
	backupCmd.Flags().IntVar(&port, "port", 5432, "Database port")
	backupCmd.Flags().StringVar(&user, "user", "", "Database user")
	backupCmd.Flags().StringVar(&password, "password", "", "Database password")
	backupCmd.Flags().StringVar(&dbname, "dbname", "", "Database name")
	backupCmd.Flags().StringVar(&output, "output", "backup.dump", "Output file name")
}
