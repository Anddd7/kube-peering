package util

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kube-peering/internal/logger"
)

func WaitElegantExit() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case sig := <-c:
			{
				logger.Z.Infof("Got %s signal. Aborting...\n", sig)
				os.Exit(1)
			}
		}
	}()
}
