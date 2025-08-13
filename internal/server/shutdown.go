package server

import (
	"os"
	"os/signal"
	"syscall"
)

// listen stop signal and close shutdounCh channel
func ListenStopSig() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
