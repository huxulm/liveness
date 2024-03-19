package main

import (
	"context"
	"path"
	"time"

	"github.com/huxulm/liveness/internal/config"
	"github.com/huxulm/liveness/internal/probex"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewApp create cobra styled command-line application
func NewApp() *cobra.Command {
	c := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			level, _ := cmd.Flags().GetUint32("loglevel")
			logrus.SetLevel(logrus.Level(level))

			dir, _ := cmd.Flags().GetString("logdir")
			logf, err := rotatelogs.New(
				path.Join(dir, "/log.%Y%m%d%H%M"),
				rotatelogs.WithLinkName(dir+"/log"),
				rotatelogs.WithRotationSize(1024*1024*50), // 50M
				rotatelogs.WithRotationTime(time.Hour),
			)
			if err != nil {
				cobra.CheckErr(err)
			}
			logrus.SetOutput(logf)

			configFile, _ := cmd.Flags().GetString("conf")
			appConf, err := config.LoadFromFile(configFile)
			cobra.CheckErr(err)
			cobra.CheckErr(probex.Run(context.TODO(), appConf))
		},
	}
	// add flags
	probex.AddProbeLoaderFlag(c)
	return c
}

func main() {
	_ = NewApp().Execute()
}
