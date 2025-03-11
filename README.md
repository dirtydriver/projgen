# RepoForge

A powerful and flexible project generator CLI tool written in Go that helps you quickly scaffold new projects from templates.

## Features

- Generate project skeletons from customizable templates
- Support for multiple project types
- Flexible parameter system via command line flags or parameter files
- Template parameter inspection
- Customizable output directory

## Installation

### From Source
```bash
# Clone the repository
git clone https://github.com/dirtydriver/RepoForge.git
cd RepoForge

# Download dependencies
make deps

# Build the binary
make build

# The binary will be available in the bin directory
```

### Via Go Install
```bash
go install github.com/dirtydriver/RepoForge@latest
```

## Development

RepoForge uses Make for common development tasks:

- `make build`: Build the binary
- `make test`: Run unit tests
- `make fmt`: Format code
- `make lint`: Run linter (requires golint)
- `make vet`: Run static analysis
- `make clean`: Clean build artifacts

## Usage

### Basic Usage

```bash
projgen --type <project-type> --name <project-name> [flags]
```

### Available Flags

- `-t, --type`: Type of project (e.g., maven, gradle, angular)
- `-n, --name`: Name of the project
- `-o, --out`: Output directory (default: current directory)
- `-p, --parameter`: Additional parameters in key=value format (can be used multiple times)
- `-f, --file`: Path to a parameters file
- `--get-template-params`: List the parameters required by the template
- `--template-dir`: Path to the template directory
- `--template-file`: Path to the template file to inspect

### Examples

1. Generate a new project:
```bash
projgen -t maven -n my-project
```

2. Generate with custom parameters:
```bash
projgen -t angular -n my-app -p version=1.0.0 -p author="John Doe"
```

3. Use a parameters file:
```bash
projgen -t gradle -n my-lib -f params.file
```

4. Check template parameters:
```bash
projgen --get-template-params -t maven --template-dir ./templates
```

## Template System

RepoForge uses a powerful templating system that allows you to create and customize project templates. Templates are stored in the `templates` directory and use the `.tmpl` extension.

### Template Structure
```
templates/
├── maven/           # Template for Maven projects
│   ├── pom.xml.tmpl
│   └── src/
├── gradle/          # Template for Gradle projects
│   ├── build.gradle.tmpl
│   └── settings.gradle.tmpl
└── angular/         # Template for Angular projects
    └── ...
```

### Parameter Files
You can create parameter files to store commonly used values. Parameter files use a simple key=value format:
```
version=1.0.0
author=John Doe
description=My awesome project
```

### Creating Custom Templates
1. Create a new directory in `templates/` for your project type
2. Add template files with the `.tmpl` extension
3. Use Go template syntax for variable substitution: `{{.variable_name}}`
4. Variables will be populated from command line parameters or parameter files

## Project Structure

```
RepoForge/
├── cmd/          # Command line interface implementation
├── filescheck/   # File system operations and checks
├── project/      # Project generation logic
├── templater/    # Template processing and rendering
├── templates/    # Project templates
├── utils/        # Utility functions
└── bin/          # Compiled binaries
```

## License

This project is licensed under the terms of the LICENSE file included in the repository.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
