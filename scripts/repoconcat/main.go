package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/unidoc/unioffice/document"
)

type Config struct {
	ImageExtensions       []string `json:"image_extensions"`
	VideoExtensions       []string `json:"video_extensions"`
	AudioExtensions       []string `json:"audio_extensions"`
	DocumentExtensions    []string `json:"document_extensions"`
	ExecutableExtensions  []string `json:"executable_extensions"`
	SettingsExtensions    []string `json:"settings_extensions"`
	AdditionalIgnoreTypes []string `json:"additional_ignore_types"`
	DefaultOutputFile     string   `json:"default_output_file"`
}

var (
	repoPath      string
	outputFile    string
	ignoreFiles   StringSlice
	ignoreTypes   StringSlice
	excludeDir    StringSlice
	ignoreSpecial bool
	includeDir    string
	config        Config
)

type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func init() {
	flag.StringVar(&repoPath, "repo_path", ".", "Path to the directory to process (ie., cloned repo). If no path is specified, defaults to the current directory.")
	flag.StringVar(&outputFile, "output_file", "", "Name for the output file. Defaults to the value in config.json.")
	flag.Var(&ignoreFiles, "ignore_files", "List of file names to ignore. Omit this argument to ignore no file names.")
	flag.Var(&ignoreTypes, "ignore_types", "List of file extensions to ignore. Defaults to the list in config.json. Omit this argument to ignore no types.")
	flag.Var(&excludeDir, "exclude_dir", "List of directory names to exclude or 'none' for no directories.")
	flag.BoolVar(&ignoreSpecial, "ignore_special", false, "Flag to ignore common settings files.")
	flag.StringVar(&includeDir, "include_dir", "", "Specific directory to include. Only contents of this directory will be documented.")
}

func loadConfig(filePath string, config *Config) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}

func shouldIgnore(item string, outputFilePath string, args Config) bool {
	itemName := filepath.Base(item)
	itemDir := filepath.Dir(item)
	fileExt := strings.ToLower(filepath.Ext(itemName))

	if filepath.Clean(item) == filepath.Clean(outputFilePath) {
		return true
	}

	if strings.HasPrefix(itemName, ".") {
		return true
	}

	// Check if the parent directory of the item is in the excludeDir list
	for _, dir := range excludeDir {
		if strings.Contains(itemDir, dir) {
			return true
		}
	}

	// if info, err := os.Stat(item); err == nil && info.IsDir() {
	// 	if contains(excludeDir, itemName) {
	// 		fmt.Println("Excluding directory:", itemName)
	// 		return true
	// 	}
	// }

	if includeDir != "" && !strings.HasPrefix(filepath.Clean(item), filepath.Clean(includeDir)) {
		return true
	}

	if info, err := os.Stat(item); err == nil && !info.IsDir() {
		if contains(ignoreFiles, itemName) || contains(ignoreTypes, fileExt) {
			return true
		}
	}

	if ignoreSpecial && contains(args.SettingsExtensions, fileExt) {
		return true
	}

	return false
}

func writeTree(dirPath string, outputFile *os.File, args Config, prefix string, isLast bool, isRoot bool) {
	if isRoot {
		fmt.Fprintf(outputFile, "%s/\n", filepath.Base(dirPath))
		isRoot = false
	}

	items, _ := ioutil.ReadDir(dirPath)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	numItems := len(items)

	for index, item := range items {
		itemPath := filepath.Join(dirPath, item.Name())

		if shouldIgnore(itemPath, outputFile.Name(), args) {
			continue
		}

		isLastItem := index == numItems-1
		newPrefix := "└── "
		childPrefix := "    "
		if !isLastItem {
			newPrefix = "├── "
			childPrefix = "│   "
		}

		fmt.Fprintf(outputFile, "%s%s%s\n", prefix, newPrefix, filepath.Base(itemPath))

		if item.IsDir() {
			nextPrefix := prefix + childPrefix
			writeTree(itemPath, outputFile, args, nextPrefix, isLastItem, false)
		}
	}
}

func writeFileContent(filePath string, outputFile *os.File, depth int) {
	indentation := strings.Repeat("  ", depth)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(outputFile, "%sError reading file: %v\n", indentation, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintf(outputFile, "%s%s\n", indentation, scanner.Text())
	}
}

func writeTreeDocx(dirPath string, doc *document.Document, args Config, outputFilePath string, prefix string, isLast bool, isRoot bool) {
	if isRoot {
		para := doc.AddParagraph()
		para.AddRun().AddText(filepath.Base(dirPath) + "/")
		isRoot = false
	}

	items, _ := ioutil.ReadDir(dirPath)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	numItems := len(items)

	for index, item := range items {
		itemPath := filepath.Join(dirPath, item.Name())

		if shouldIgnore(itemPath, outputFilePath, args) {
			continue
		}

		isLastItem := index == numItems-1
		newPrefix := "└── "
		childPrefix := "    "
		if !isLastItem {
			newPrefix = "├── "
			childPrefix = "│   "
		}

		para := doc.AddParagraph()
		para.AddRun().AddText(prefix + newPrefix + filepath.Base(itemPath))

		if item.IsDir() {
			nextPrefix := prefix + childPrefix
			writeTreeDocx(itemPath, doc, args, outputFilePath, nextPrefix, isLastItem, false)
		}
	}
}

func writeFileContentDocx(filePath string, doc *document.Document) {
	file, err := os.Open(filePath)
	if err != nil {
		para := doc.AddParagraph()
		para.AddRun().AddText(fmt.Sprintf("Error reading file: %v", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		para := doc.AddParagraph()
		para.AddRun().AddText(scanner.Text())
	}
}

func writeFileContentsInOrder(dirPath string, outputFile *os.File, args Config, depth int) {
	items, _ := ioutil.ReadDir(dirPath)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	for _, item := range items {
		itemPath := filepath.Join(dirPath, item.Name())
		relativePath, _ := filepath.Rel(repoPath, itemPath)

		if shouldIgnore(itemPath, outputFile.Name(), args) {
			continue
		}

		if item.IsDir() {
			writeFileContentsInOrder(itemPath, outputFile, args, depth+1)
		} else {
			fmt.Fprintf(outputFile, "%s[File Begins] %s\n", strings.Repeat("  ", depth), relativePath)
			writeFileContent(itemPath, outputFile, depth)
			fmt.Fprintf(outputFile, "\n%s[File Ends] %s\n\n", strings.Repeat("  ", depth), relativePath)
		}
	}
}

func writeFileContentsInOrderDocx(dirPath string, doc *document.Document, args Config, depth int) {
	items, _ := ioutil.ReadDir(dirPath)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	for _, item := range items {
		itemPath := filepath.Join(dirPath, item.Name())
		relativePath, _ := filepath.Rel(repoPath, itemPath)

		if shouldIgnore(itemPath, outputFile, args) {
			continue
		}

		if item.IsDir() {
			writeFileContentsInOrderDocx(itemPath, doc, args, depth+1)
		} else {
			doc.AddParagraph().AddRun().AddText(fmt.Sprintf("[File Begins] %s", relativePath))
			writeFileContentDocx(itemPath, doc)
			doc.AddParagraph().AddRun().AddText(fmt.Sprintf("[File Ends] %s", relativePath))
		}
	}
}

func main() {
	flag.Parse()

	// Debug print to check the parsed values
	fmt.Println("Exclude directories:", excludeDir)

	err := loadConfig("config.json", &config)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	if outputFile == "" {
		outputFile = config.DefaultOutputFile
	}

	if ignoreTypes == nil {
		ignoreTypes = append(config.ImageExtensions, config.VideoExtensions...)
		ignoreTypes = append(ignoreTypes, config.AudioExtensions...)
		ignoreTypes = append(ignoreTypes, config.DocumentExtensions...)
		ignoreTypes = append(ignoreTypes, config.ExecutableExtensions...)
		ignoreTypes = append(ignoreTypes, config.AdditionalIgnoreTypes...)
	}

	if !filepath.IsAbs(repoPath) {
		repoPath, _ = filepath.Abs(repoPath)
	}

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("Error: The specified directory does not exist or is not accessible: %s\n", repoPath)
		return
	}

	if strings.HasSuffix(outputFile, ".docx") {
		doc := document.New()

		doc.AddParagraph().AddRun().AddText("Repository Documentation")
		doc.AddParagraph().AddRun().AddText("This document provides a comprehensive overview of the repository's structure and contents. The first section, titled 'Directory/File Tree', displays the repository's hierarchy in a tree format. In this section, directories and files are listed using tree branches to indicate their structure and relationships. Following the tree representation, the 'File Content' section details the contents of each file in the repository. Each file's content is introduced with a '[File Begins]' marker followed by the file's relative path, and the content is displayed verbatim. The end of each file's content is marked with a '[File Ends]' marker. This format ensures a clear and orderly presentation of both the structure and the detailed contents of the repository.")

		doc.AddParagraph().AddRun().AddText("Directory/File Tree Begins -->")
		writeTreeDocx(repoPath, doc, config, outputFile, "", true, true)
		doc.AddParagraph().AddRun().AddText("<-- Directory/File Tree Ends")

		doc.AddParagraph().AddRun().AddText("File Content Begins -->")
		writeFileContentsInOrderDocx(repoPath, doc, config, 0)
		doc.AddParagraph().AddRun().AddText("<-- File Content Ends")

		doc.SaveToFile(outputFile)
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer file.Close()

		fmt.Fprintln(file, "Repository Documentation")
		fmt.Fprintln(file, "This document provides a comprehensive overview of the repository's structure and contents. The first section, titled 'Directory/File Tree', displays the repository's hierarchy in a tree format. In this section, directories and files are listed using tree branches to indicate their structure and relationships. Following the tree representation, the 'File Content' section details the contents of each file in the repository. Each file's content is introduced with a '[File Begins]' marker followed by the file's relative path, and the content is displayed verbatim. The end of each file's content is marked with a '[File Ends]' marker. This format ensures a clear and orderly presentation of both the structure and the detailed contents of the repository.")

		fmt.Fprintln(file, "Directory/File Tree Begins -->")
		writeTree(repoPath, file, config, "", true, true)
		fmt.Fprintln(file, "<-- Directory/File Tree Ends")

		fmt.Fprintln(file, "File Content Begins -->")
		writeFileContentsInOrder(repoPath, file, config, 0)
		fmt.Fprintln(file, "<-- File Content Ends")
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
