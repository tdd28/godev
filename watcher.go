package main

import (
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	*fsnotify.Watcher

	handler  func(fsnotify.Event)
	skipDirs []string
}

func (w *Watcher) Watch(root string) error {
	defer w.Close()

	go w.watch()

	if err := w.walkDir(root); err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-sig

	return nil
}

func (w *Watcher) watch() {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Println(err)
		case ev, ok := <-w.Events:
			if !ok {
				return
			}

			switch ev.Op {
			case fsnotify.Create:
				if fi, err := os.Stat(ev.Name); err == nil && fi.IsDir() {
					w.walkDir(ev.Name)
				}
			}

			w.handler(ev)
		}
	}
}

func (w *Watcher) walkDir(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			for _, skipDir := range w.skipDirs {
				if d.Name() == skipDir {
					return filepath.SkipDir
				}
			}
			if err := w.Add(path); err != nil {
				log.Println(err)
			}
		}

		return nil
	})
}
