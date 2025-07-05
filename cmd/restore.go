package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a PostgreSQL database from a backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting restore operation...")

		cmdArgs := []string{
			"-h", host,
			"-p", strconv.Itoa(port),
			"-U", user,
			"-d", restoreDb,
			"--clean",
			"--if-exists",
			"--verbose",
			inputFile,
		}

		restoreCmd := exec.Command("pg_restore", cmdArgs...)
		restoreCmd.Env = append(os.Environ(), "PGPASSWORD="+password)

		restoreCmd.Stdout = os.Stdout
		restoreCmd.Stderr = os.Stderr

		err := restoreCmd.Run()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}

		fmt.Println("Restore complete.")
	},
}

var (
	inputFile string
	restoreDb string
)

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringVar(&host, "host", "localhost", "Database host")
	restoreCmd.Flags().IntVar(&port, "port", 5432, "Database port")
	restoreCmd.Flags().StringVar(&user, "user", "", "Database user")
	restoreCmd.Flags().StringVar(&password, "password", "", "Database password")
	restoreCmd.Flags().StringVar(&inputFile, "input", "", "Path to backup file (.dump)")
	restoreCmd.Flags().StringVar(&restoreDb, "dbname", "", "Target database to restore into")

	restoreCmd.MarkFlagRequired("input")
	restoreCmd.MarkFlagRequired("dbname")
}
