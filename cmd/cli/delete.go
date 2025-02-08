package cli

import (
	"fmt"
	"log"

	"github.com/CinematicCow/lumora/internal/core"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "rm <key>",
	Short: "Delete a key-value pair",
	Long:  "Delete a key-value pair",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		db, err := core.Open(dataDir)
		if err != nil {
			log.Fatalf("Open failed: %v", err)
		}
		defer db.Close()

		err = db.Delete(key)
		if err != nil {
			log.Fatalf("Get failed: %v", err)
		}
		fmt.Printf("Key: %s deleted\n", key)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
