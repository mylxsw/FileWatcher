package main

import (
	"bytes"
	"flag"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os/exec"
	"strings"
)

var watch_path = flag.String("path", "", "监控的目录")
var cmd = flag.String("cmd", "", "要执行的命令")
var suffix = flag.String("suffix", ".mmd", "监控文件后缀")

func main() {
	flag.Parse()
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
					if strings.HasSuffix(event.Name, *suffix) {
						log.Println("modified file: ", event.Name)
						//cmd_to_exec := "/usr/local/bin/mermaid -w 1200 -t" + watch_path + "/mermaid.css -g " + watch_path + "/config.json -o " + watch_path + "/output " + event.Name
						//cmd_to_exec := fmt.Sprintf(*cmd, event.Name)
						cmd_to_exec := strings.Replace(*cmd, "[filename]", event.Name, -1)
						cmd_to_exec = strings.Replace(cmd_to_exec, "[path]", *watch_path, -1)
						log.Println("cmd to execute: ", cmd_to_exec)

						cmd := exec.Command("sh", "-c", cmd_to_exec)
						var out_stdout bytes.Buffer
						var out_stderr bytes.Buffer
						cmd.Stderr = &out_stderr
						cmd.Stdout = &out_stdout
						err := cmd.Run()
						if err != nil {
							log.Println("command execute failed: ", err)
						}
						log.Println(out_stdout.String())
						if out_stderr.Len() > 0 {
							log.Println("Error: ", out_stderr.String())
						}
					}

				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*watch_path)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("watching directory: " + *watch_path)

	<-done
}
