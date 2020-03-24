package cmdutil

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/XSAM/go-hybrid/log"
	"github.com/XSAM/go-hybrid/metadata"
)

func VersionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			log.BgLogger().Info("application info", zap.Any("info", metadata.AppInfo()))
		},
	}
	return &cmd
}
