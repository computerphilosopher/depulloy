package main

import (
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

type Daemon interface {
	Run()
}

type ZKDaemon struct {
	zkConn *zk.Conn
	zkRoot string

	versionWatcher Watcher
	copyWatcher    Watcher
	runWatcher     Watcher
	cmdWatcher     Watcher
}

func NewZKDaemon(zkServersRaw, zkRoot string, timeout time.Duration) (*ZKDaemon, error) {
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

	return &ZKDaemon{
		zkConn:         conn,
		zkRoot:         zkRoot,
		versionWatcher: versionWatcher,
		copyWatcher:    copyWatcher,
		runWatcher:     runWatcher,
		cmdWatcher:     cmdWatcher,
	}, nil
}

func (daemon *ZKDaemon) Run() {
}
