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
- [Contributing](#-contributing)
  - [Quick Start for Contributors](#quick-start-for-contributors)
  - [Development Setup](#development-setup)
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

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Install development dependencies
go mod download

# Install pre-commit hooks
pre-commit install

# Run linting
golangci-lint run

# Run tests
go test ./...
```

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
