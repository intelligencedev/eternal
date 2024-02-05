package documents

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	golangssh "golang.org/x/crypto/ssh"
)

type Document struct {
	PageContent string
	Metadata    map[string]string
}

type GitLoader struct {
	RepoPath           string
	CloneURL           string
	Branch             string
	PrivateKeyPath     string
	FileFilter         func(string) bool
	InsecureSkipVerify bool
}

func NewGitLoader(repoPath, cloneURL, branch, privateKeyPath string, fileFilter func(string) bool, insecureSkipVerify bool) *GitLoader {
	return &GitLoader{RepoPath: repoPath, CloneURL: cloneURL, Branch: branch, PrivateKeyPath: privateKeyPath, FileFilter: fileFilter, InsecureSkipVerify: insecureSkipVerify}
}

// Load loads the documents from the Git repository specified by the GitLoader.
// It returns a slice of Document and an error if any.
// If the repository does not exist at the specified path and a clone URL is provided,
// it clones the repository using the provided authentication options.
// If the repository already exists, it opens the repository at the specified path.
// If a branch is specified, it checks out the branch.
// It then walks through the repository files, reads the content of each file,
// and creates a Document object for each file with the corresponding metadata.
// The resulting documents are returned as a slice.
// If any error occurs during the process, it is returned.
func (gl *GitLoader) Load() ([]Document, error) {
	var repo *gogit.Repository
	var err error

	if _, err = os.Stat(gl.RepoPath); os.IsNotExist(err) && gl.CloneURL != "" {
		sshKey, _ := os.ReadFile(gl.PrivateKeyPath)
		signer, _ := golangssh.ParsePrivateKey(sshKey)
		auth := &gitssh.PublicKeys{User: "git", Signer: signer}
		if gl.InsecureSkipVerify {
			auth.HostKeyCallback = golangssh.InsecureIgnoreHostKey()
		}
		repo, err = gogit.PlainClone(gl.RepoPath, false, &gogit.CloneOptions{URL: gl.CloneURL, Auth: auth})
		if err != nil {
			return nil, err
		}
	} else {
		repo, err = gogit.PlainOpen(gl.RepoPath)
		if err != nil {
			return nil, err
		}
	}

	if gl.Branch != "" {
		w, err := repo.Worktree()
		if err != nil {
			return nil, err
		}
		err = w.Checkout(&gogit.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(gl.Branch)})
		if err != nil {
			return nil, err
		}
	}

	var docs []Document

	err = filepath.Walk(gl.RepoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if gl.FileFilter != nil && !gl.FileFilter(path) {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		textContent := string(content)
		relFilePath, _ := filepath.Rel(gl.RepoPath, path)
		fileType := filepath.Ext(info.Name())

		metadata := map[string]string{
			"source":    relFilePath,
			"file_path": relFilePath,
			"file_name": info.Name(),
			"file_type": fileType,
		}

		doc := Document{PageContent: textContent, Metadata: metadata}
		docs = append(docs, doc)

		return nil
	})

	if err != nil {
		fmt.Printf("Error reading files: %s\n", err)
	}

	return docs, nil
}
