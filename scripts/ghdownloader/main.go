package main

import (
	"eternal/pkg/ghdownloader"
	"fmt"
)

func main() {
	err := ghdownloader.DownloadAndExtractRepo("comfyanonymous", "ComfyUI", "", "/Users/arturoaquino/Downloads")
	if err != nil {
		fmt.Printf("Error downloading repository: %v\n", err)
		return
	}
	fmt.Println("Repository downloaded successfully!")
}
