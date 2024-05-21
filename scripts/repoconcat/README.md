# Repoconcat

Concatenates an entire git code repo into one text file for ML RAG workflows.

```
# Example command with flags for Eternal repo
$ go run . -repo_path /Users/$USER/eternal -ignore_files eternal -ignore_files LICENSE -ignore_files Makefile -ignore_files public/js/drawflow.min.js -ignore_files public/js/htmx.min.js -ignore_files public/js/package-lock.json -ignore_files public/js/package.json -exclude_dir index -exclude_dir examples -exclude_dir services -exclude_dir tmp -exclude_dir reference -exclude_dir scripts -exclude_dir public/css/drawflow -exclude_dir public/css/halfmoon -exclude_dir public/fonts -exclude_dir public/img -exclude_dir public/js/bootstrap -exclude_dir public/js/highlight -exclude_dir public/js/node_modules -exclude_dir pkg/llm/local/gguf -exclude_dir pkg/llm/local/bin -exclude_dir pkg/llm/local/gguf -exclude_dir pkg/sd/sdcpp -ignore_types .jpg -ignore_types .png -ignore_types .mod -ignore_types .sum -ignore_types .log -ignore_types .md -ignore_types .txt
```