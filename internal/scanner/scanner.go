package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type FileInfo struct {
	Path            string
	Type            FileType
	HasCopyright    bool
	CopyrightYear   int
	CopyrightNotice string
	LineNumber      int
}

type FileType struct {
	Name         string
	Extensions   []string
	CommentStart string
	CommentEnd   string
	LineComment  string
}

var supportedFileTypes = []FileType{
	{
		Name:        "Go",
		Extensions:  []string{".go"},
		LineComment: "//",
	},
	{
		Name:        "Python",
		Extensions:  []string{".py"},
		LineComment: "#",
	},
	{
		Name:        "JavaScript/TypeScript",
		Extensions:  []string{".js", ".ts", ".jsx", ".tsx"},
		LineComment: "//",
	},
	{
		Name:        "Java",
		Extensions:  []string{".java"},
		LineComment: "//",
	},
	{
		Name:        "C/C++",
		Extensions:  []string{".c", ".cpp", ".cc", ".cxx", ".h", ".hpp"},
		LineComment: "//",
	},
	{
		Name:        "Shell",
		Extensions:  []string{".sh", ".bash"},
		LineComment: "#",
	},
}

// Scanner handles file scanning and copyright detection
type Scanner struct {
	copyrightRegex *regexp.Regexp
	yearRegex      *regexp.Regexp
	verbose        bool
}

// NewScanner creates a new scanner instance
func NewScanner(verbose bool) *Scanner {
	return &Scanner{
		copyrightRegex: regexp.MustCompile(`(?i)copyright\s*(\(c\))?\s*(\d{4})?`),
		yearRegex:      regexp.MustCompile(`\b(19|20)\d{2}\b`),
		verbose:        verbose,
	}
}

// ScanFile analyzes a file for copyright information
func (s *Scanner) ScanFile(filePath string) (*FileInfo, error) {
	fileType := s.detectFileType(filePath)
	if fileType == nil {
		return nil, fmt.Errorf("unsupported file type: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	info := &FileInfo{
		Path: filePath,
		Type: *fileType,
	}

	s.analyzeCopyright(string(content), info)
	return info, nil
}

// ScanFiles processes multiple files
func (s *Scanner) ScanFiles(filePaths []string) ([]*FileInfo, error) {
	var results []*FileInfo
	var errors []error

	for _, path := range filePaths {
		// Handle directories
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			dirFiles, err := s.scanDirectory(path)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			results = append(results, dirFiles...)
		} else {
			// Handle single file
			info, err := s.ScanFile(path)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			results = append(results, info)
		}
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("encountered %d errors during scanning", len(errors))
	}

	return results, nil
}

// scanDirectory recursively scans a directory
func (s *Scanner) scanDirectory(dirPath string) ([]*FileInfo, error) {
	var results []*FileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if s.detectFileType(path) != nil {
			fileInfo, scanErr := s.ScanFile(path)
			if scanErr != nil {
				fmt.Printf("Warning: failed to scan %s: %v\n", path, scanErr)
			} else {
				results = append(results, fileInfo)
			}
		}

		return nil
	})

	return results, err
}

// detectFileType determines the file type based on extension
func (s *Scanner) detectFileType(filePath string) *FileType {
	ext := strings.ToLower(filepath.Ext(filePath))

	for _, fileType := range supportedFileTypes {
		for _, supportedExt := range fileType.Extensions {
			if ext == supportedExt {
				return &fileType
			}
		}
	}

	return nil
}

// analyzeCopyright looks for copyright notices in file content
func (s *Scanner) analyzeCopyright(content string, info *FileInfo) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	// Only check first 20 lines for copyright
	for scanner.Scan() && lineNum < 20 {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Remove comment markers
		cleanLine := s.removeCommentMarkers(line, info.Type)

		// Debug: Print what we're checking (remove this later)
		if s.verbose && strings.Contains(strings.ToLower(cleanLine), "copyright") {
			fmt.Printf("DEBUG: Found 'copyright' in line %d: '%s' -> cleaned: '%s'\n", lineNum, line, cleanLine)
		}

		// Check if this line contains copyright
		if s.copyrightRegex.MatchString(cleanLine) {
			info.HasCopyright = true
			info.CopyrightNotice = line
			info.LineNumber = lineNum

			// Try to extract year
			if years := s.yearRegex.FindAllString(cleanLine, -1); len(years) > 0 {
				// Use the last year found (most recent)
				if year := parseInt(years[len(years)-1]); year > 0 {
					info.CopyrightYear = year
				}
			}
			break
		}
	}
}

// removeCommentMarkers strips comment syntax from a line
func (s *Scanner) removeCommentMarkers(line string, fileType FileType) string {
	cleaned := line

	// Remove line comment marker
	if fileType.LineComment != "" {
		cleaned = strings.TrimPrefix(cleaned, fileType.LineComment)
		cleaned = strings.TrimSpace(cleaned)
	}

	// Remove block comment markers if present
	if fileType.CommentStart != "" {
		cleaned = strings.TrimPrefix(cleaned, fileType.CommentStart)
	}
	if fileType.CommentEnd != "" {
		cleaned = strings.TrimSuffix(cleaned, fileType.CommentEnd)
	}

	return strings.TrimSpace(cleaned)
}

// parseInt safely converts string to int
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

// IsOutdated checks if the copyright year is outdated
func (info *FileInfo) IsOutdated() bool {
	if !info.HasCopyright || info.CopyrightYear == 0 {
		return false
	}
	return info.CopyrightYear < time.Now().Year()
}

// NeedsUpdate returns true if the file needs copyright updates
func (info *FileInfo) NeedsUpdate() bool {
	return !info.HasCopyright || info.IsOutdated()
}
