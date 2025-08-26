![kafta logo](img/kafta.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/electric-saw/kafta)](https://goreportcard.com/report/github.com/electric-saw/kafta)
[![License](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE.txt)
[![Release](https://img.shields.io/github/release/electric-saw/kafta.svg)](https://github.com/electric-saw/kafta/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/electric-saw/kafta)](https://golang.org/)

A modern, **non-JVM** command-line interface for managing **Apache Kafka** clusters. Inspired by `kubectl`, Kafta provides a simple and efficient way to manage topics, brokers, consumer groups, and more across multiple Kafka clusters.

## Table of Contents

- [Why Kafta?](#-why-kafta)
- [Installation](#-installation)
  - [Prerequisites](#prerequisites)
  - [Quick Install](#quick-install)
  - [macOS Installation Guide](#-macos-installation-guide)
  - [Linux Installation](#-linux-installation)
  - [Windows Installation](#-windows-installation)
- [Troubleshooting Installation](#-troubleshooting-installation)
  - [Common Issues](#common-issues)
  - [Diagnostic Script](#-diagnostic-script)
- [Configuration](#️-configuration)
  - [Initial Setup](#initial-setup)
  - [Configuration with Authentication](#configuration-with-authentication)
  - [Sample Configuration File](#sample-configuration-file)
- [Usage Examples](#-usage-examples)
  - [Cluster Management](#cluster-management)
  - [Topic Management](#topic-management)
  - [Consumer Group Management](#consumer-group-management)
  - [Message Production & Consumption](#message-production--consumption)
  - [Schema Registry Operations](#schema-registry-operations)
  - [Output Formats](#output-formats)
- [Smart Features](#-smart-features)
  - [Intelligent Suggestions](#intelligent-suggestions)
  - [Command Completion](#command-completion)
- [Advanced Features](#️-advanced-features)
  - [Environment Variables](#environment-variables)
  - [Configuration Templates](#configuration-templates)
- [Development](#️-development)
  - [Building from Source](#building-from-source)
  - [Running Tests](#running-tests)
  - [Development Setup](#development-setup)
    - [macOS Development Setup](#-macos-development-setup)
    - [Linux Development Setup](#-linux-development-setup)
    - [Windows Development Setup](#-windows-development-setup)
  - [Development Environment Verification](#-development-environment-verification)
  - [Quick Start for Contributors](#️-quick-start-for-contributors)
  - [Docker Development Environment](#-docker-development-environment)
  - [Development Workflow](#-development-workflow)
  - [Code Quality Tools](#-code-quality-tools)
  - [Testing Guidelines](#-testing-guidelines)
  - [Performance Profiling](#-performance-profiling)
- [Contributing](#-contributing)
- [Roadmap](#-roadmap)
- [Known Issues](#-known-issues)
- [License](#-license)
- [Acknowledgments](#-acknowledgments)
- [Support & Community](#-support--community)
- [Star History](#-star-history)

## 🚀 Why Kafta?

- **🔥 No Java Required**: Binary executable, no JVM setup needed
- **🎯 kubectl-inspired**: Familiar interface for Kubernetes users  
- **⚡ Fast & Lightweight**: Written in Go for speed and efficiency
- **🔧 Multi-cluster**: Easy switching between different Kafka environments
- **📊 Rich Output**: Table, JSON, and YAML output formats
- **🤖 Smart Suggestions**: Intelligent command corrections and hints
- **🔐 Security First**: Built-in support for SASL authentication and TLS

---

## 📦 Installation

### Prerequisites

- **Go 1.18+** (for installation from source)
- **Apache Kafka cluster** (for usage)

### Quick Install

#### For Go >= 1.18 (Recommended)
```bash
go install github.com/electric-saw/kafta/cmd/kafta@latest
```

#### For Go < 1.18 (Legacy)
```bash
go get -u github.com/electric-saw/kafta
```

### 🍎 macOS Installation Guide

If you're on macOS and encountering PATH issues, follow these steps:

#### 1. Install Kafta
```bash
go install github.com/electric-saw/kafta/cmd/kafta@latest
```

#### 2. Configure PATH (Choose your shell)

**For zsh (default on modern macOS):**
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

**For bash:**
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
source ~/.bash_profile
```

#### 3. Verify Installation
```bash
kafta --help
```

### 🐧 Linux Installation

```bash
# Install
go install github.com/electric-saw/kafta/cmd/kafta@latest

# Add to PATH (usually automatic, but if needed)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### 🪟 Windows Installation

```powershell
# Install
go install github.com/electric-saw/kafta/cmd/kafta@latest

# The binary will be in %GOPATH%\bin or %USERPROFILE%\go\bin
# Add to PATH if necessary through System Properties > Environment Variables
```

---

## 🔧 Troubleshooting Installation

### Common Issues

#### ❌ "kafta: command not found"

**Diagnosis:**
```bash
# Check if kafta exists
ls -la $(go env GOPATH)/bin/kafta

# Check current PATH
echo $PATH
```

**Solution:**
```bash
# Temporary fix
export PATH=$PATH:$(go env GOPATH)/bin

# Permanent fix (add to shell config)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

#### ❌ Permission denied

```bash
# Fix permissions
chmod +x $(go env GOPATH)/bin/kafta
```

#### ❌ Go not found

Install Go from [golang.org](https://golang.org/dl/) or:

```bash
# macOS with Homebrew
brew install go

# Ubuntu/Debian
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang
```

### 🩺 Diagnostic Script

Run this to diagnose installation issues:

```bash
#!/bin/bash
echo "🔍 Kafta Installation Diagnostics"
echo "================================="

echo "Go version: $(go version 2>/dev/null || echo 'Go not found')"
echo "GOPATH: $(go env GOPATH 2>/dev/null || echo 'Go not configured')"

GOPATH=$(go env GOPATH 2>/dev/null)
if [ -f "$GOPATH/bin/kafta" ]; then
    echo "✅ Kafta found at: $GOPATH/bin/kafta"
else
    echo "❌ Kafta not found in GOPATH/bin"
fi

if command -v kafta >/dev/null 2>&1; then
    echo "✅ Kafta available in PATH"
    kafta --version 2>/dev/null || echo "⚠️ Kafta found but with errors"
else
    echo "❌ Kafta not in PATH"
    echo "💡 Run: export PATH=\$PATH:$(go env GOPATH)/bin"
fi
```

---

## ⚙️ Configuration

Kafta stores configurations in `~/.kafta/config` as a YAML file to support multiple Kafka clusters.

### Initial Setup

```bash
# Configure your first Kafka cluster
kafta config set-context my-cluster \
  --bootstrap-servers "localhost:9092" \
  --schema-registry "http://localhost:8081"

# Use the configured context
kafta config use-context my-cluster

# Verify current context
kafta config current-context
```

### Configuration with Authentication

```bash
# Configure cluster with SASL authentication
kafta config set-context production \
  --bootstrap-servers "broker1:9092,broker2:9092,broker3:9092" \
  --schema-registry "https://schema-registry.prod.com" \
  --use-sasl true \
  --sasl-algorithm "sha512" \
  --user "your-username" \
  --password "your-password"
```

### Sample Configuration File

```yaml
# ~/.kafta/config
contexts:
  development:
    bootstrap_servers: "localhost:9092"
    schema_registry: "http://localhost:8081"
    use_sasl: false
  
  production:
    bootstrap_servers: "prod-broker1:9092,prod-broker2:9092"
    schema_registry: "https://schema-registry.prod.com"
    use_sasl: true
    sasl_algorithm: "sha512"
    user: "kafta-user"
    password: "secure-password"

current_context: "development"
```

---

## 🎯 Usage Examples

### Cluster Management

```bash
# List all configured contexts
kafta config get-contexts

# Switch between clusters
kafta config use-context production

# Describe current cluster
kafta cluster describe
```

### Topic Management

```bash
# List all topics
kafta topic list

# Create a new topic
kafta topic create my-topic \
  --partitions 3 \
  --replication-factor 2

# Describe topic details
kafta topic describe my-topic

# Delete a topic
kafta topic delete my-topic

# List topic configurations
kafta topic list-configs my-topic
```

### Consumer Group Management

```bash
# List all consumer groups
kafta consumer list

# Check consumer lag
kafta consumer lag my-consumer-group

# Describe consumer group
kafta consumer describe my-consumer-group

# Delete consumer group
kafta consumer delete my-consumer-group
```

### Message Production & Consumption

```bash
# Produce messages to a topic
kafta producer --topic my-topic
# Type your messages and press Enter
# Use Ctrl+C to exit

# Consume messages from a topic
kafta console consumer --topic my-topic

# Consume with specific consumer group
kafta console consumer --topic my-topic --group my-group

# Consume with verbose output
kafta console consumer --topic my-topic --verbose
```

### Schema Registry Operations

```bash
# List all schema subjects
kafta schema subjects-list

# Get latest schema for a subject
kafta schema get my-topic-value

# Get specific version of schema
kafta schema get my-topic-value --version 2

# List all versions for a subject
kafta schema versions my-topic-value

# Compare schema versions (shows diff)
kafta schema diff my-topic-value --from-version 1 --to-version 2
```

### Output Formats

```bash
# Default table output
kafta topic list

# JSON output
kafta topic list --output json

# YAML output  
kafta topic list --output yaml

# Save output to file
kafta topic list --output json > topics.json
```

---

## 🎪 Smart Features

### Intelligent Suggestions

Kafta provides helpful suggestions when you make typos:

```bash
$ kafta clustr describe
Command 'clustr' not found. Did you mean 'cluster'?

$ kafta topic craete my-topic
Command 'craete' not found. Did you mean 'create'?
```

### Command Completion

Enable bash/zsh completion:

```bash
# For bash
kafta completion bash > /etc/bash_completion.d/kafta

# For zsh  
kafta completion zsh > "${fpath[1]}/_kafta"

# For fish
kafta completion fish > ~/.config/fish/completions/kafta.fish
```

---

## 🏗️ Advanced Features

### Environment Variables

Override configuration with environment variables:

```bash
export KAFTA_BOOTSTRAP_SERVERS="localhost:9092"
export KAFTA_SCHEMA_REGISTRY="http://localhost:8081"
export KAFTA_SASL_USERNAME="my-user"
export KAFTA_SASL_PASSWORD="my-password"

kafta cluster describe  # Uses environment variables
```

### Configuration Templates

Create reusable configuration templates:

```bash
# Create template for development environments
kafta config create-template dev-template \
  --bootstrap-servers "{{.host}}:9092" \
  --schema-registry "http://{{.host}}:8081"

# Apply template
kafta config apply-template dev-template \
  --set host=localhost \
  --context my-dev-cluster
```

---

## 🛠️ Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/electric-saw/kafta.git
cd kafta

# Build binary
go build -o kafta cmd/kafta/main.go

# Install locally
go install ./cmd/kafta
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires Kafka)
go test -tags=integration ./...
```

### Development Setup

#### 🍎 macOS Development Setup

```bash
# 1. Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Install Go (if not already installed)
brew install go

# 3. Install development tools
brew install golangci-lint pre-commit

# 4. Clone and setup project
git clone https://github.com/electric-saw/kafta.git
cd kafta

# 5. Install Go dependencies
go mod download

# 6. Setup pre-commit hooks
pre-commit install

# 7. Run initial checks
golangci-lint run
go test ./...

# 8. Verify installation
echo "✅ Development environment ready!"
kafta --version || echo "Build the project first: go build -o kafta cmd/kafta/main.go"
```

#### 🐧 Linux Development Setup

**Ubuntu/Debian:**
```bash
# 1. Update package manager
sudo apt update

# 2. Install Go (if not already installed)
sudo apt install golang-go

# 3. Install Python and pip (for pre-commit)
sudo apt install python3 python3-pip

# 4. Install pre-commit
pip3 install pre-commit

# 5. Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

# 6. Add Go bin to PATH
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# 7. Clone and setup project
git clone https://github.com/electric-saw/kafta.git
cd kafta

# 8. Install Go dependencies
go mod download

# 9. Setup pre-commit hooks
pre-commit install

# 10. Run initial checks
golangci-lint run
go test ./...
```

**CentOS/RHEL/Fedora:**
```bash
# 1. Install Go (if not already installed)
sudo dnf install golang  # Fedora
# OR
sudo yum install golang  # CentOS/RHEL

# 2. Install Python and pip
sudo dnf install python3 python3-pip  # Fedora
# OR
sudo yum install python3 python3-pip  # CentOS/RHEL

# 3. Install pre-commit
pip3 install --user pre-commit

# 4. Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

# 5. Add to PATH
echo 'export PATH=$PATH:$(go env GOPATH)/bin:$HOME/.local/bin' >> ~/.bashrc
source ~/.bashrc

# 6. Continue with project setup (same as Ubuntu)
git clone https://github.com/electric-saw/kafta.git
cd kafta
go mod download
pre-commit install
golangci-lint run
go test ./...
```

#### 🪟 Windows Development Setup

**PowerShell (Run as Administrator):**
```powershell
# 1. Install Chocolatey (package manager for Windows)
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 2. Install Go (if not already installed)
choco install golang

# 3. Install Python (for pre-commit)
choco install python

# 4. Install Git (if not already installed)
choco install git

# 5. Refresh environment variables
refreshenv

# 6. Install pre-commit
pip install pre-commit

# 7. Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 8. Clone and setup project
git clone https://github.com/electric-saw/kafta.git
cd kafta

# 9. Install Go dependencies
go mod download

# 10. Setup pre-commit hooks
pre-commit install

# 11. Run initial checks
golangci-lint run
go test ./...
```

**Alternative Windows Setup (without Chocolatey):**
```powershell
# 1. Download and install Go manually from https://golang.org/dl/
# 2. Download and install Python from https://python.org/downloads/
# 3. Download and install Git from https://git-scm.com/download/win

# 4. Install pre-commit via pip
pip install pre-commit

# 5. Install golangci-lint via Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 6. Ensure Go bin is in PATH (usually automatic, but check)
# Add %USERPROFILE%\goin to your PATH environment variable if needed

# 7. Continue with project setup
git clone https://github.com/electric-saw/kafta.git
cd kafta
go mod download
pre-commit install
golangci-lint run
go test ./...
```

### 🩺 Development Environment Verification

Run this script to verify your development setup:

```bash
#!/bin/bash
echo "🔍 Kafta Development Environment Check"
echo "====================================="

echo ""
echo "📦 Go Environment:"
go version 2>/dev/null || echo "❌ Go not installed"
echo "GOPATH: $(go env GOPATH 2>/dev/null)"
echo "GOROOT: $(go env GOROOT 2>/dev/null)"

echo ""
echo "🔧 Development Tools:"
golangci-lint --version 2>/dev/null || echo "❌ golangci-lint not installed"
pre-commit --version 2>/dev/null || echo "❌ pre-commit not installed"
git --version 2>/dev/null || echo "❌ git not installed"

echo ""
echo "📁 Project Setup:"
[ -f "go.mod" ] && echo "✅ go.mod found" || echo "❌ go.mod not found"
[ -f ".golangci.yml" ] && echo "✅ .golangci.yml found" || echo "❌ .golangci.yml not found"
[ -f ".pre-commit-config.yaml" ] && echo "✅ .pre-commit-config.yaml found" || echo "❌ .pre-commit-config.yaml not found"

echo ""
echo "🧪 Running Quick Tests:"
echo "• go mod tidy: $(go mod tidy 2>&1 && echo "✅ OK" || echo "❌ Failed")"
echo "• go build: $(go build cmd/kafta/main.go 2>&1 && echo "✅ OK" || echo "❌ Failed")"
echo "• golangci-lint: $(golangci-lint run --timeout=60s 2>&1 >/dev/null && echo "✅ OK" || echo "⚠️ Issues found")"

echo ""
if command -v kafta >/dev/null 2>&1; then
    echo "✅ kafta command available: $(kafta --version 2>/dev/null || echo "version command not available")"
else
    echo "ℹ️  kafta not in PATH (build and install: go install ./cmd/kafta)"
fi
```

### 🏃‍♂️ Quick Start for Contributors

```bash
# One-liner setup for macOS
brew install go golangci-lint pre-commit && git clone https://github.com/electric-saw/kafta.git && cd kafta && go mod download && pre-commit install

# One-liner setup for Ubuntu/Debian
sudo apt update && sudo apt install golang-go python3-pip && pip3 install pre-commit && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Verify everything works
go test ./... && golangci-lint run && echo "🎉 Ready to contribute!"
```

### 🐳 Docker Development Environment

If you prefer to use Docker for development:

```bash
# Build development Docker image
docker build -t kafta-dev -f Dockerfile.dev .

# Run development container
docker run -it --rm -v $(pwd):/workspace kafta-dev

# Inside container, run development commands
go test ./...
golangci-lint run
```

### 📋 Development Workflow

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/kafta.git
   cd kafta
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Setup development tools**:
   ```bash
   pre-commit install
   ```
5. **Create a feature branch**:
   ```bash
   git checkout -b feature/my-awesome-feature
   ```
6. **Make your changes** and test:
   ```bash
   go test ./...
   golangci-lint run
   ```
7. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add awesome feature"
   ```
8. **Push and create Pull Request**:
   ```bash
   git push origin feature/my-awesome-feature
   ```

### 🔍 Code Quality Tools

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run linter with auto-fix
golangci-lint run --fix

# Run pre-commit on all files
pre-commit run --all-files

# Update dependencies
go mod tidy
go mod vendor  # if using vendor

# Security scan
go list -json -m all | nancy sleuth  # requires nancy to be installed
```

### 🧪 Testing Guidelines

```bash
# Run specific test
go test ./pkg/cmd/...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmark tests
go test -bench=. ./...
```

### 📈 Performance Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# Build with profiling enabled
go build -tags=profile ./cmd/kafta
```

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📋 Roadmap

### ✅ Completed Features
- [x] Basic topic management (CRUD)
- [x] Consumer group management  
- [x] Multi-cluster configuration
- [x] SASL authentication support
- [x] Schema Registry integration
- [x] Smart command suggestions
- [x] Multiple output formats

### 🚧 Work In Progress
- [ ] Advanced schema management (evolution, compatibility)
- [ ] KSQL support and management
- [ ] Real-time topic data tailing
- [ ] Performance benchmarking tools
- [ ] Monitoring and alerting
- [ ] ACL management

### 🎯 Planned Features
- [ ] Kafka Connect management
- [ ] Transaction support
- [ ] Kafka Streams integration
- [ ] Web UI dashboard
- [ ] Plugin system
- [ ] Docker/Kubernetes integration

---

## 🐛 Known Issues

- Schema registry operations may timeout with large schemas
- Consumer lag calculation might be inaccurate during rebalancing
- Some SASL mechanisms are not yet supported

See [Issues](https://github.com/electric-saw/kafta/issues) for the latest bugs and feature requests.

---

## 📄 License

This project is licensed under the AGPL-3.0 License - see the [LICENSE.txt](LICENSE.txt) file for details.

---

## 🙏 Acknowledgments

- Inspired by [kubectl](https://kubernetes.io/docs/reference/kubectl/) for the command structure
- Built with [Sarama](https://github.com/IBM/sarama) Kafka library
- Uses [Cobra](https://github.com/spf13/cobra) for CLI framework

---

## 📞 Support & Community

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/electric-saw/kafta/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/electric-saw/kafta/discussions)  
- 📖 **Documentation**: [Wiki](https://github.com/electric-saw/kafta/wiki)
- 💬 **Community**: [Slack Channel](https://join.slack.com/t/kafta-community/shared_invite/...)

---

## ⭐ Star History

If you find Kafta useful, please consider giving it a star! ⭐

[![Star History Chart](https://api.star-history.com/svg?repos=electric-saw/kafta&type=Date)](https://star-history.com/#electric-saw/kafta&Date)

---

**Made with ❤️ by the Kafta team**
