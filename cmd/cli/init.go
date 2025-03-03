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

		dataDir = filepath.Join(os.Getenv("HOME"), ".config", "lumora", dbName)

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatalf("Failed to create data directory: %v", err)
		}

		setDefault, _ := cmd.Flags().GetBool("default")
		cfg.AddDB(dbName, dataDir)
		if setDefault || cfg.DefaultDB == "" {
			cfg.SetDefaultDB(dbName)
		}

		if err := cfg.Save(); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}

		_, err = core.Open(dataDir)
		if err != nil {
			log.Fatalf("Initialization failed: %v", err)
		}
		fmt.Printf("Database initialized at %s/\n", dataDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().Bool("default", false, "Set this DB as default")
}
