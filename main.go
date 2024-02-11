package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

func main() {

	zkAddr := flag.String("zk-addr", "localhost:2181", "comma seperated list of zookeeper address")
	zkPath := flag.String("zk-path", "/depulloy/dev/example", "znode path of app")
	timeoutRaw := flag.String("timeout", "20s", "znode path of app")

	timeout, err := time.ParseDuration(*timeoutRaw)
	if err != nil {
		panic(fmt.Sprintf("invalid timeout flag %s: %s", *timeoutRaw, err.Error()))
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	daemon, err := NewZKDaemon(*zkAddr, *zkPath, logger, timeout)
	if err != nil {
		panic(fmt.Sprintf("cannot get new daemon: %s", timeoutRaw, err.Error()))
	}

	daemon.Run()
}
