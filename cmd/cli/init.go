package cli

import (
	"fmt"
	"log"

	"github.com/CinematicCow/lumora/internal/core"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize lumora",
	Long:  "Initialize lumora in the specified directory",
	Run: func(cmd *cobra.Command, args []string) {
		if dataDir == "" {
			dataDir, _ = cmd.Flags().GetString("data-dir")
		}
		_, err := core.Open(dataDir)
		if err != nil {
			log.Fatalf("Initialization failed: %v", err)
		}
		fmt.Printf("Database initialized at %s/\n", dataDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
