package client

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

var DefaultIgnorePatterns = []gitignore.Pattern{
	gitignore.ParsePattern(MetadataFilename, nil), // Knowledge Metadata file
	gitignore.ParsePattern("~$*", nil),            // MS Office temp files
	gitignore.ParsePattern("$*", nil),             // Likely hidden/tempfiles
}

func isIgnored(ignore gitignore.Matcher, path string) bool {
	return ignore.Match(strings.Split(path, string(filepath.Separator)), false)
}

func readIgnoreFile(path string) ([]gitignore.Pattern, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to checkout ignore file %q: %w", path, err)
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("ignore file %q is a directory", path)
	}

	var ps []gitignore.Pattern
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open ignore file %q: %w", path, err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, "#") && len(strings.TrimSpace(s)) > 0 {
			ps = append(ps, gitignore.ParsePattern(s, nil))
		}
	}

	return ps, nil
}
