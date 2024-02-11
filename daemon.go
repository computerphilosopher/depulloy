package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

type Daemon struct {
	zkConn *zk.Conn
	zkRoot string

	versionWatcher Watcher
	copyWatcher    Watcher
	runWatcher     Watcher
	cmdWatcher     Watcher
}

func (daemon *Daemon) getPath(terminalNode string) string {
	return fmt.Sprintf("%s/%s", daemon.zkRoot, terminalNode)
}

func NewDaemon(zkServersRaw, zkRoot string, timeout time.Duration) (*Daemon, error) {
	zkServers := strings.Split(zkServersRaw, ",")
	conn, _, err := zk.Connect(zkServers, timeout)
	if err != nil {
		return nil, err
	}

	versionWatcher, err := NewZKWatcher(conn, zkRoot+"/version")
	if err != nil {
		return nil, err
	}

	copyWatcher, err := NewZKWatcher(conn, zkRoot+"/run")
	if err != nil {
		return nil, err
	}

	runWatcher, err := NewZKWatcher(conn, zkRoot+"/run")
	if err != nil {
		return nil, err
	}

	cmdWatcher, err := NewZKWatcher(conn, zkRoot+"/cmd")
	if err != nil {
		return nil, err
	}

	return &Daemon{
		zkConn:         conn,
		zkRoot:         zkRoot,
		versionWatcher: versionWatcher,
		copyWatcher:    copyWatcher,
		runWatcher:     runWatcher,
		cmdWatcher:     cmdWatcher,
	}, nil
}

func (daemon *Daemon) Run() {
}
