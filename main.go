package main

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	args := []string{"run"}
	args = append(args, os.Args[1:]...)

	c := newCmd("go", args...)
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
	defer c.Stop()

	timer := time.AfterFunc(math.MaxInt64, func() { c.Restart() })

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	watcher := &Watcher{
		Watcher:  w,
		skipDirs: []string{".git"},
		handler: func(e fsnotify.Event) {
			switch e.Op {
			case fsnotify.Chmod:
			default:
				timer.Reset(time.Second)
			}
		},
	}

	watcher.Watch(root)
}
