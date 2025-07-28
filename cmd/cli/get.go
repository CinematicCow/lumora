package cli

import (
	"errors"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return WithDDK(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		dataDir := cmd.Context().Value(DDK).(string)

		db, err := core.Open(dataDir)
		if err != nil {
			log.Fatalf("Open failed: %v", err)
		}
		defer db.Close()

		value, err := db.Get(key)
		if err != nil {
			if errors.Is(err, core.ErrKeyNotFound) {
				fmt.Printf("key %q not found in %q\n", key, dataDir)
			}
			log.Fatalf("Get failed: %v", err)
		}
		fmt.Printf("Key: %s | Value: %s\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("name", "n", "", "database name to query")
}
