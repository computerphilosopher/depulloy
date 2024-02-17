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
	logger *slog.Logger

	versionWatcher *ZKWatcher
	copyWatcher    *ZKWatcher
	envWatcher     *ZKWatcher
	runWatcher     *ZKWatcher
}

func NewZKDaemon(zkServersRaw, zkRoot string, logger *slog.Logger, timeout time.Duration) (*ZKDaemon, error) {
	zkServers := strings.Split(zkServersRaw, ",")
	conn, _, err := zk.Connect(zkServers, timeout)
	if err != nil {
		return nil, err
	}

	copyWatcher, err := NewZKWatcher(NewZKNode(conn, zkRoot+"/copy"), logger)
	if err != nil {
		return nil, err
	}

	envWatcher, err := NewZKWatcher(NewZKNode(conn, zkRoot+"/env"), logger)
	if err != nil {
		return nil, err
	}

	runWatcher, err := NewZKWatcher(NewZKNode(conn, zkRoot+"/run"), logger)
	if err != nil {
		return nil, err
	}

	ret := &ZKDaemon{
		zkConn:      conn,
		zkRoot:      zkRoot,
		logger:      logger,
		copyWatcher: copyWatcher,
		envWatcher:  envWatcher,
		runWatcher:  runWatcher,
	}
	versionWatcher, err := NewZKWatcher(NewZKNode(conn, zkRoot+"/version"), logger)
	if err != nil {
		return nil, err
	}

	ret.versionWatcher = versionWatcher
	return ret, nil
}

func (daemon *ZKDaemon) Run() {
	go daemon.copyWatcher.Watch()
	go daemon.envWatcher.Watch()
	go daemon.runWatcher.Watch()
	go daemon.versionWatcher.WatchWithRetry()
}

func (daemon *ZKDaemon) Deploy() error {
	return nil
}
