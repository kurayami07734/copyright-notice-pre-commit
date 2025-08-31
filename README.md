# Copyright Notice Pre-commit

This pre-commit hook automatically adds a copyright notice to source code files during a commit. It also updates the year in the copyright notice, ensuring it is always current.

## Features

- Adds a copyright notice at the top of files if not present.
- Automatically updates the year in the copyright notice to the current year.
- Supports multiple programming languages (Go, Python, JavaScript, TypeScript, Java, C/C++, Shell)
- Configurable company name and notice format
- Integration with pre-commit framework
- Dry-run mode for previewing changes

## Installation

### Build from source
```bash
git clone https://github.com/kurayami07734/copyright-notice-pre-commit
cd copyright-notice-pre-commit
go build -o copyright cmd/copyright/main.go
```

### Install globally
```bash
go install ./cmd/copyright
```

## Usage

```
Copyright Notice CLI

Usage:
  copyright [command] [flags] [files...]

Commands:
  check     Check files for copyright notices
  fix       Add/update copyright notices
  version   Show version information

Examples:
  copyright check src/
  copyright fix --auto-fix --company "Acme Inc" **/*.go
  copyright check --config .copyright.yaml .
```

### Commands

#### `check` - Validate copyright notices
```bash
# Check files in current directory
copyright check .

# Check specific files with verbose output
copyright check --verbose src/main.go src/utils.go

# Check with custom company name
copyright check --company "Your Company" src/

# Use custom config file
copyright check --config .copyright.yaml src/
```

**Flags:**
- `--config` - Path to configuration file
- `--company` - Company name (overrides config)
- `--verbose` - Show detailed output

#### `fix` - Add/update copyright notices
```bash
# Automatically fix all issues
copyright fix --auto-fix src/

# Preview changes without making them
copyright fix --dry-run src/

# Fix with custom company name
copyright fix --auto-fix --company "Acme Inc" **/*.go

# Fix using configuration file
copyright fix --auto-fix --config .copyright.yaml .
```

**Flags:**
- `--auto-fix` - Automatically add/update copyright notices
- `--dry-run` - Show what would be changed without making changes
- `--config` - Path to configuration file
- `--company` - Company name (overrides config)

#### `version` - Show version information
```bash
copyright version
```

## Configuration

Create a `.copyright.yaml` file in your project root:

```yaml
company_name: "Your Company Name"
notice_format: "Copyright (C) $year $company_name. All rights reserved."
auto_fix: false

file_patterns:
  - "*.go"
  - "*.py"
  - "*.js"
  - "*.ts"
  - "*.java"
  - "*.cpp"
  - "*.c"
  - "*.h"
  - "*.sh"

exclude_patterns:
  - "vendor/"
  - "node_modules/"
  - ".git/"
  - "*.pb.go"
  - "*_generated.go"
  - "*.min.js"
```

### Template Variables

The `notice_format` supports these variables:
- `$year` - Current year
- `$current_year` - Current year (alias)
- `$company_name` - Company name from config

### Example Notice Formats

```yaml
# Standard format
notice_format: "Copyright (C) $year $company_name. All rights reserved."

# Simple format
notice_format: "Copyright $year $company_name"

# With license info
notice_format: "Copyright (C) $year $company_name. Licensed under MIT."
```

## Supported File Types

| Language | Extensions | Comment Style |
|----------|------------|---------------|
| Go | `.go` | `//` |
| Python | `.py` | `#` |
| JavaScript/TypeScript | `.js`, `.ts`, `.jsx`, `.tsx` | `//` |
| Java | `.java` | `//` |
| C/C++ | `.c`, `.cpp`, `.cc`, `.cxx`, `.h`, `.hpp` | `//` |
| Shell | `.sh`, `.bash` | `#` |

## Pre-commit Integration

### Using pre-commit framework

Add to your `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: copyright-check
        name: Check Copyright Notices
        entry: copyright check
        language: golang
        files: \.(go|py|js|ts|java|cpp|c|h|sh)$
        pass_filenames: true
```

### Manual git hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash
set -e

# Get list of staged files
files=$(git diff --cached --name-only --diff-filter=AM)

if [ -z "$files" ]; then
    exit 0
fi

# Run copyright check
./copyright check $files

if [ $? -ne 0 ]; then
    echo "Copyright check failed. Run 'copyright fix --auto-fix' to automatically fix issues."
    exit 1
fi
```

```bash
chmod +x .git/hooks/pre-commit
```

## Examples

### Basic workflow
```bash
# 1. Check current status
copyright check --verbose .

# 2. Fix issues automatically
copyright fix --auto-fix --company "My Company" .

# 3. Verify fixes
copyright check .
```

### Integration with CI/CD
```bash
# In your CI pipeline
copyright check .
if [ $? -ne 0 ]; then
  echo "Copyright notices are missing or outdated"
  exit 1
fi
```

## Troubleshooting

### Common Issues

**Files not being processed:**
- Check if file extension is supported
- Verify file is not in exclude patterns
- Use `--verbose` flag to see what files are being scanned

**Copyright not detected:**
- Ensure copyright notice is in the first 20 lines
- Check that proper comment syntax is used for the file type
- Use `--verbose` to see detection details

**Wrong company name:**
- Use `--company` flag to override
- Update `.copyright.yaml` configuration file 