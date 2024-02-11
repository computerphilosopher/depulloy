package main

import (
	"fmt"

	"github.com/go-zookeeper/zk"
)

// TODO: MySQL, Postgres Watcher
type Watcher interface {
	Watch()
}

type ZKWatcher struct {
	conn    *zk.Conn
	zkPath  string
	content []byte
	watcher <-chan zk.Event
}

func NewZKWatcher(conn *zk.Conn, zkPath string) (*ZKWatcher, error) {
	content, _, watcher, err := conn.GetW(zkPath)
	if err != nil {
		return nil, err
	}
	return &ZKWatcher{
		conn:    conn,
		zkPath:  zkPath,
		content: content,
		watcher: watcher,
	}, nil
}

func (watcher *ZKWatcher) Watch() {
	for {
		select {
		case event := <-watcher.watcher:
			err := watcher.setNewWatcher(event)
			if err != nil {
				continue
			}
		}
	}
}

func (watcher *ZKWatcher) setNewWatcher(event zk.Event) error {
	newContent, _, newWatcher, err := watcher.conn.GetW(watcher.zkPath)
	if err != nil {
		watcher.content = nil
		watcher.watcher = nil
		return err
	}

	watcher.content = newContent
	watcher.watcher = newWatcher

	if event.Type != zk.EventNodeDataChanged {
		return fmt.Errorf("event is %s", event.Type.String())
	}

	return nil
}