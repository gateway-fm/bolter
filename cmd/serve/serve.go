package serve

import (
	"github.com/spf13/cobra"

	"bolter/internal"
)

// Cmd is
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fire",
		Short: "Start Bolter!!!",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.LoadBolter(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
