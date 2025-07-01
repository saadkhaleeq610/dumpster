package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	utils "github.com/saadkhaleeq610/dumpster/utility"
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

		logFile, err := os.OpenFile("dumpster.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()

		logger := log.New(logFile, "", log.LstdFlags)
		logger.Println("🔁 Starting backup for DB:", dbname)
		start := time.Now()

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

		err = dumpCmd.Run()
		if err != nil {
			log.Fatalf("Backup failed: %v", err)
		}

		fmt.Printf("Backup complete: %s\n", output)

		compressed := output + ".gz"
		err = utils.CompressFile(output, compressed)
		if err != nil {
			log.Fatalf("Compression failed: %v", err)
		}

		fmt.Printf("Compressed: %s\n", compressed)

		os.Remove(output)

		duration := time.Since(start)
		logger.Printf("Backup complete: %s | Size: %s | Duration: %s\n", compressed, fileSize(compressed), duration)
	},
}

func fileSize(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "unknown"
	}
	return fmt.Sprintf("%.2f MB", float64(info.Size())/1024.0/1024.0)
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
