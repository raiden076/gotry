package git

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const commitMessage = `✨ Let's try something new

🤖 Created with gotry (https://github.com/raiden076/gotry)`

func Init(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitialCommit(path string) error {
	// Create .gitkeep to have something to commit
	gitkeep := filepath.Join(path, ".gitkeep")
	if err := os.WriteFile(gitkeep, []byte{}, 0644); err != nil {
		return err
	}

	// git add .
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = path
	if err := addCmd.Run(); err != nil {
		return err
	}

	// git commit
	commitCmd := exec.Command("git", "commit", "-m", commitMessage)
	commitCmd.Dir = path
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	return commitCmd.Run()
}

func Clone(repoURL, destPath string) error {
	cmd := exec.Command("git", "clone", repoURL, destPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type RepoInfo struct {
	Host string
	User string
	Repo string
}

func ParseGitURL(rawURL string) (*RepoInfo, error) {
	// Handle SSH format: git@github.com:user/repo.git
	sshRegex := regexp.MustCompile(`^git@([^:]+):([^/]+)/(.+?)(?:\.git)?$`)
	if matches := sshRegex.FindStringSubmatch(rawURL); matches != nil {
		return &RepoInfo{
			Host: matches[1],
			User: matches[2],
			Repo: matches[3],
		}, nil
	}

	// Handle HTTPS format
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid repository URL: %s", rawURL)
	}

	repo := pathParts[1]
	repo = strings.TrimSuffix(repo, ".git")

	return &RepoInfo{
		Host: parsed.Host,
		User: pathParts[0],
		Repo: repo,
	}, nil
}

func (r *RepoInfo) DirectoryName() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s-%s-%s", today, r.User, r.Repo)
}

func IsGitURL(s string) bool {
	return strings.HasPrefix(s, "git@") ||
		strings.HasPrefix(s, "https://github.com") ||
		strings.HasPrefix(s, "https://gitlab.com") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") && strings.Contains(s, ".git")
}
