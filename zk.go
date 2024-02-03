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

	version        []byte
	versionWatcher <-chan zk.Event

	copy        []byte
	copyWatcher <-chan zk.Event

	run        []byte
	runWatcher <-chan zk.Event

	cmd        []byte
	cmdWatcher <-chan zk.Event
}

var eventStrByType map[zk.EventType]string = map[zk.EventType]string{
	zk.EventNodeDataChanged: "NODE_DATA_CHANGED",
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

	version, _, versionWatcher, err := conn.GetW(zkRoot + "/version")
	if err != nil {
		return nil, err
	}

	copy, _, copyWatcher, err := conn.GetW(zkRoot + "/copy")
	if err != nil {
		return nil, err
	}

	run, _, runWatcher, err := conn.GetW(zkRoot + "/run")
	if err != nil {
		return nil, err
	}

	cmd, _, cmdWatcher, err := conn.GetW(zkRoot + "/cmd")
	if err != nil {
		return nil, err
	}

	return &Daemon{
		zkConn:         conn,
		zkRoot:         zkRoot,
		version:        version,
		versionWatcher: versionWatcher,
		copy:           copy,
		copyWatcher:    copyWatcher,
		run:            run,
		runWatcher:     runWatcher,
		cmd:            cmd,
		cmdWatcher:     cmdWatcher,
	}, nil
}

func (daemon *Daemon) Run() {
	for {
		select {
		case event := <-daemon.copyWatcher:
			daemon.receiveEventAndSetNewWatcher(event, daemon.getPath("/copy"), daemon.copy, daemon.copyWatcher)
			daemon.deploy()
		case <-daemon.copyWatcher:

		}
	}
}

func (daemon *Daemon) receiveEventAndSetNewWatcher(event zk.Event, zkPath string, content []byte, watcher <-chan zk.Event) error {

	newContent, _, newWatcher, err := daemon.zkConn.GetW(zkPath)
	if err != nil {
		content = nil
		watcher = nil
		return err
	}

	content = newContent
	watcher = newWatcher

	if event.Type != zk.EventNodeDataChanged {
		eventStr, exist := eventStrByType[event.Type]
		if !exist {
			return fmt.Errorf("unknown event type: %d", event.Type)
		}
		return fmt.Errorf("event is %s", eventStr)
	}
	return nil
}

func (daemon *Daemon) deploy() {
}
