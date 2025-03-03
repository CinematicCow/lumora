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
	Use:   "init",
	Short: "Initialize lumora",
	Long:  "Initialize lumora in the specified directory",
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.InitConfig()
		if err != nil {
			log.Fatalf("failed to initialize config: %v", err)
		}

		if dataDir == "" {
			dataDir = cfg.DBPaths[cfg.DefaultDB]
			if dataDir == "" {
				dataDir = filepath.Join(os.Getenv("HOME"), ".config", "lumora")
			}
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
}
