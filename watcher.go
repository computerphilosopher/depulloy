package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-zookeeper/zk"
)

type ZKWatcher struct {
	znode     *ZKNode
	eventChan <-chan zk.Event
	logger    *slog.Logger
}

func NewZKWatcher(znode *ZKNode, logger *slog.Logger) (*ZKWatcher, error) {
	watcher, err := znode.GetW()
	if err != nil {
		return nil, err
	}
	return &ZKWatcher{
		znode:     znode,
		eventChan: watcher,
		logger:    logger,
	}, nil
}

func (watcher *ZKWatcher) handleEvent(event zk.Event) error {
	if event.Type != zk.EventNodeDataChanged {
		return fmt.Errorf("event is %s", event.Type.String())
	}

	eventChan, err := watcher.znode.GetW()
	if err != nil {
		return err
	}
	watcher.eventChan = eventChan

	return nil
}

func (watcher *ZKWatcher) Watch() {
	for {
		select {
		case event := <-watcher.eventChan:
			err := watcher.handleEvent(event)
			if err != nil {
				watcher.logger.Error(err.Error())
			}
		}
	}
}

func (watcher *ZKWatcher) WatchWithRetry() {
	for {
		select {
		case event := <-watcher.eventChan:
			for {
				err := watcher.handleEvent(event)
				if err == nil {
					break
				}
				watcher.logger.Error(err.Error())
				time.Sleep(time.Second * 5)
			}
		}
	}
}
