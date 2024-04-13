# Repoconcat

Concatenates an entire git code repo into one text file for ML RAG workflows.

```
# Example command with flags for Eternal repo
$ go run . -repo_path /Users/art/Documents/eternal-v1 -ignore_files eternal -ignore_files LICENSE -exclude_dir reference -exclude_dir scripts -exclude_dir public -exclude_dir pkg/llm/local/gguf -exclude_dir pkg/llm/local/bin -exclude_dir pkg/llm/local/gguf -exclude_dir pkg/sd/sdcpp -ignore_types .jpg -ignore_types .png -ignore_types .mod -ignore_types .sum -ignore_types .log -ignore_types .md -ignore_types .txt
```