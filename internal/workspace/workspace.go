package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Directory struct {
	Name      string
	Path      string
	ModTime   time.Time
	DatePart  string
	NamePart  string
}

func List(basePath string) ([]Directory, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Directory{}, nil
		}
		return nil, err
	}

	var dirs []Directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		datePart, namePart := parseDirectoryName(name)

		dirs = append(dirs, Directory{
			Name:     name,
			Path:     filepath.Join(basePath, name),
			ModTime:  info.ModTime(),
			DatePart: datePart,
			NamePart: namePart,
		})
	}

	// Sort by modification time, most recent first
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].ModTime.After(dirs[j].ModTime)
	})

	return dirs, nil
}

func parseDirectoryName(name string) (datePart, namePart string) {
	// Expected format: YYYY-MM-DD-name
	if len(name) >= 11 && name[4] == '-' && name[7] == '-' && name[10] == '-' {
		return name[:10], name[11:]
	}
	return "", name
}

func Create(basePath, name string) (string, error) {
	today := time.Now().Format("2006-01-02")
	dirName := fmt.Sprintf("%s-%s", today, sanitizeName(name))
	fullPath := filepath.Join(basePath, dirName)

	// Handle collisions
	finalPath := fullPath
	counter := 2
	for {
		if _, err := os.Stat(finalPath); os.IsNotExist(err) {
			break
		}
		finalPath = fmt.Sprintf("%s-%d", fullPath, counter)
		counter++
	}

	if err := os.MkdirAll(finalPath, 0755); err != nil {
		return "", err
	}

	return finalPath, nil
}

func sanitizeName(name string) string {
	// Replace spaces with hyphens, lowercase
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	return name
}

func Delete(paths []string) error {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func RelativeTime(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "now"
	case duration < time.Hour:
		mins := int(duration.Minutes())
		return fmt.Sprintf("%dm", mins)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	default:
		weeks := int(duration.Hours() / 24 / 7)
		return fmt.Sprintf("%dw", weeks)
	}
}
