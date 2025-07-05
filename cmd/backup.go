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
	Short: "Creates a timestamped PostgreSQL backup and compresses it",
	Run: func(cmd *cobra.Command, args []string) {
		logFile, err := os.OpenFile("dumpster.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()

		logger := log.New(logFile, "", log.LstdFlags)
		logger.Println("Starting backup for DB:", dbname)
		start := time.Now()

		fmt.Println("Starting PostgreSQL backup...")

		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("backup_%s.dump", timestamp)

		cmdArgs := []string{
			"-h", host,
			"-p", strconv.Itoa(port),
			"-U", user,
			"-F", "c",
			"-f", filename,
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

		fmt.Printf("Backup complete: %s\n", filename)

		compressed := filename + ".gz"
		err = utils.CompressFile(filename, compressed)
		if err != nil {
			log.Fatalf("Compression failed: %v", err)
		}

		fmt.Printf("Compressed: %s\n", compressed)

		if uploadToS3 {
			fmt.Println("Uploading to S3...")
			err = utils.UploadToS3(compressed, bucket, region)
			if err != nil {
				log.Fatalf("S3 upload failed: %v", err)
			}
			fmt.Println("S3 upload complete.")
		}

		os.Remove(filename)

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
	host       string
	port       int
	user       string
	password   string
	dbname     string
	bucket     string
	region     string
	uploadToS3 bool
)

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVar(&host, "host", "localhost", "Database host")
	backupCmd.Flags().IntVar(&port, "port", 5432, "Database port")
	backupCmd.Flags().StringVar(&user, "user", "", "Database user")
	backupCmd.Flags().StringVar(&password, "password", "", "Database password")
	backupCmd.Flags().StringVar(&dbname, "dbname", "", "Database name")
	backupCmd.Flags().BoolVar(&uploadToS3, "s3", false, "Upload backup to AWS S3")
	backupCmd.Flags().StringVar(&bucket, "bucket", "", "AWS S3 bucket name")
	backupCmd.Flags().StringVar(&region, "region", "us-east-1", "AWS region")

}
