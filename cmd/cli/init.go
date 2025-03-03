package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/CinematicCow/lumora/internal/config"
	"github.com/CinematicCow/lumora/internal/core"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [db-name]",
	Short: "Initialize lumora",
	Long:  "Initialize lumora with default database",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.InitConfig()
		if err != nil {
			log.Fatalf("failed to initialize config: %v", err)
		}

		dbName := "default"
		if len(args) > 0 {
			dbName = args[0]
		}

		if dataDir == "" {
			dataDir = filepath.Join(os.Getenv("HOME"), ".config", "lumora", dbName)
		}

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}

		cfg.AddDB(dbName, dataDir)
		cfg.SetDefaultDB(dbName)

		_, err = core.Open(dataDir)
		if err != nil {
			log.Fatalf("Initialization failed: %v", err)
		}
		fmt.Printf("Database initialized at %s/\n", dataDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
