package cli

import (
	"context"
	"fmt"

	"github.com/CinematicCow/lumora/internal/config"
	"github.com/spf13/cobra"
)

// context key
type ck string

// Data Dir Key
const DDK ck = "dataDir"

func WithDDK(cmd *cobra.Command) error {

	cfg, err := config.InitConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	dbName, _ := cmd.Flags().GetString("name")
	if dbName == "" {
		dbName = cfg.DefaultDB
	}

	if dbName == "" {
		return fmt.Errorf("no database specified. Either:\n1. Use --db-name flag\n2. Initialize new database: lumora init <name>\n")
	}

	dbPath, exists := cfg.GetDBPath(dbName)
	if !exists {
		return fmt.Errorf("database %q not found.\nAvailable options:\n- Initialize it: lumora init %s\n", dbName, dbName)
	}

	cmd.SetContext(context.WithValue(cmd.Context(), DDK, dbPath))
	return nil
}
