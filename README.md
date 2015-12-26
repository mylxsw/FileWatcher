# FileWatcher

文件变动监控，在目录下的文件发生变动的时候，自动触发配置的命令，比如编译等

下列命令监控目录`/Users/mylxsw/Documents/Works/note`下的所有`*.mmd`文件，如果文件发生修改，则执行`mermaid`命令，为该文件创建新的流程图，输出到`/Users/mylxsw/Documents/Works/note/output`目录下。

    go run watcher.go -cmd="mermaid -w 1200 -t [path]/mermaid.css -g [path]/config.json -o [path]/output [filename]" -path=/Users/mylxsw/Documents/Works/note

- `[path]` 工作目录，会被替换成`-path`参数指定的路径
- `[filename]` 发生修改的文件名（带绝对路径信息）

