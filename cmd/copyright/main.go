package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kurayami07734/copyright-notice-pre-commit/internal/config"
	"github.com/kurayami07734/copyright-notice-pre-commit/internal/scanner"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "check":
		runCheck(args)
	case "fix":
		runFix(args)
	case "version":
		runVersion()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showUsage()
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("Copyright Notice CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  copyright [command] [flags] [files...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  check     Check files for copyright notices")
	fmt.Println("  fix       Add/update copyright notices")
	fmt.Println("  version   Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  copyright check src/")
	fmt.Println("  copyright fix --auto-fix --company \"Acme Inc\" **/*.go")
	fmt.Println("  copyright check --config .copyright.yaml .")
}

func runCheck(args []string) {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	configFile := fs.String("config", "", "Path to config file")
	company := fs.String("company", "", "Company name")
	verbose := fs.Bool("verbose", false, "Verbose output")

	fs.Parse(args)
	files := fs.Args()

	if len(files) == 0 {
		fmt.Println("No files specified")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override with command line flags
	cfg.OverrideFromFlags(*company, "", false)

	// Create scanner
	s := scanner.NewScanner(*verbose)

	// Scan files
	results, err := s.ScanFiles(files)
	if err != nil {
		fmt.Printf("Error scanning files: %v\n", err)
	}

	// Report results
	var missingCount, outdatedCount int
	for _, result := range results {
		if !result.HasCopyright {
			missingCount++
			fmt.Printf("MISSING: %s\n", result.Path)
		} else if result.IsOutdated() {
			outdatedCount++
			fmt.Printf("OUTDATED: %s (year: %d)\n", result.Path, result.CopyrightYear)
		} else {
			fmt.Printf("OK: %s\n", result.Path)
		}
	}

	fmt.Printf("Scanned %d files: %d missing copyright, %d outdated\n", len(results), missingCount, outdatedCount)

	if missingCount > 0 || outdatedCount > 0 {
		fmt.Println("Run with 'fix --auto-fix' to automatically fix issues")
		os.Exit(1)
	}
}

func runFix(args []string) {
	fs := flag.NewFlagSet("fix", flag.ExitOnError)
	configFile := fs.String("config", "", "Path to config file")
	company := fs.String("company", "", "Company name")
	autoFix := fs.Bool("auto-fix", false, "Automatically fix issues")
	dryRun := fs.Bool("dry-run", false, "Show what would be changed without making changes")

	fs.Parse(args)
	files := fs.Args()

	fmt.Printf("Fixing copyright notices...\n")
	fmt.Printf("Config: %s, Company: %s, Auto-fix: %t, Dry-run: %t\n", *configFile, *company, *autoFix, *dryRun)
	fmt.Printf("Files: %v\n", files)

	// TODO: Implement actual fix logic
	fmt.Println("Fix functionality not yet implemented")
}

func runVersion() {
	fmt.Printf("copyright-cli version %s\n", version)
}
