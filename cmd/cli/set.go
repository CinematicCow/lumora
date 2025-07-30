package cli

import (
	"fmt"
	"log"

	"github.com/CinematicCow/lumora/internal/core"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Store a key-value pair",
	Long:  "Store a key-value pair",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return WithDDK(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		dataDir := cmd.Context().Value(DDK).(string)

		db, err := core.Open(dataDir)
		if err != nil {
			log.Fatalf("Open failed: %v", err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("Close failed: %v", err)
			}
		}()

		if err := db.Put(key, []byte(value)); err != nil {
			log.Fatalf("Set failed: %v", err)
		}
		fmt.Printf("Key: %s | Value: %s\n", key, value)
	},
}

func init() {
	rootCmd.AddCommand(putCmd)
	putCmd.Flags().StringP("name", "n", "", "database name to use")
}
