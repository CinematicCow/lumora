package cli

import (
	"fmt"
	"log"

	"github.com/CinematicCow/lumora/internal/core"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a value via key",
	Long:  "Get a value via key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		db, err := core.Open(dataDir)
		if err != nil {
			log.Fatalf("Open failed: %v", err)
		}
		defer db.Close()

		value, err := db.Get(key)
		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
		fmt.Printf("Key: %s | Value: %s\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
