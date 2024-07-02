# Concatenates text files into one file

```
$ go run main.go -paths="/path1,/path2" -types=".txt,.go" -ignore="_test" -output="result.txt" -recursive=true
```

go run main.go -paths="/Users/arturoaquino/Documents/eternal,/Users/arturoaquino/Documents/eternal/pkg/llm/local,/Users/arturoaquino/Documents/eternal/public/templates" -types=".html,.css.js,.go" -ignore="_test" -output="result.txt" -recursive=false