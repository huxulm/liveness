package log

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func AddLogrusFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("logdir", "logdir", ".log", "log directory")
	cmd.Flags().Uint32P("loglevel", "loglevel", uint32(log.DebugLevel), "log directory")
}
