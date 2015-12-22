package main

import (
	"bytes"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os/exec"
	"strings"
)

func main() {
	watch_path := "/tmp/test/note"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println(event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
					if strings.HasSuffix(event.Name, ".mmd") {
						cmd_to_exec := "/usr/local/bin/mermaid -t" + watch_path + "/mermaid.forest.css -g " + watch_path + "/config.json -o " + watch_path + "/output " + event.Name
						log.Println("cmd to execute: ", cmd_to_exec)

						cmd := exec.Command("sh", "-c", cmd_to_exec)
						var out bytes.Buffer
						var out2 bytes.Buffer
						cmd.Stderr = &out2
						cmd.Stdout = &out
						err := cmd.Run()
						if err != nil {
							log.Println("command execute failed: ", err)
						}
						log.Println(out.String(), out2.String())
					}

				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watch_path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
