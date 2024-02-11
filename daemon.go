package main

import (
	"log/slog"
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

func NewZKDaemon(zkServersRaw, zkRoot string, logger *slog.Logger, timeout time.Duration) (*ZKDaemon, error) {
	zkServers := strings.Split(zkServersRaw, ",")
	conn, _, err := zk.Connect(zkServers, timeout)
	if err != nil {
		return nil, err
	}

	copyWatcher, err := NewZKWatcher(conn, zkRoot+"/copy", logger)
	if err != nil {
		return nil, err
	}

	runWatcher, err := NewZKWatcher(conn, zkRoot+"/run", logger)
	if err != nil {
		return nil, err
	}

	cmdWatcher, err := NewZKWatcher(conn, zkRoot+"/cmd", logger)
	if err != nil {
		return nil, err
	}

	ret := &ZKDaemon{
		zkConn:      conn,
		zkRoot:      zkRoot,
		copyWatcher: copyWatcher,
		runWatcher:  runWatcher,
		cmdWatcher:  cmdWatcher,
	}
	versionWatcher, err := NewZKWatcher(conn, zkRoot+"/version", logger)
	if err != nil {
		return nil, err
	}

	ret.versionWatcher = versionWatcher
	return ret, nil
}

func (daemon *ZKDaemon) Run() {
	go daemon.versionWatcher.Watch()
	go daemon.copyWatcher.Watch()
	go daemon.runWatcher.Watch()
	go daemon.cmdWatcher.Watch()
}

func (daemon *ZKDaemon) Deploy() error {
	return nil
}
