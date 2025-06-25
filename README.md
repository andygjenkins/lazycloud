# LazyCloud 🚀

AWS Terminal UI for the Modern Developer - inspired by lazydocker, lazygit, and k9s.

## Quick Start

### 🧪 Development (LocalStack)

Get started in 30 seconds with a complete test environment:

```bash
# Complete setup with test Lambda functions, S3 buckets, ECS resources
make dev-setup

# Run LazyCloud with test data
make dev-run
```

**What you get:**
- 3 test Lambda functions ready to explore
- S3 buckets with sample data
- ECS cluster with task definitions
- No AWS credentials needed

### ☁️ Production (AWS)

Connect to your real AWS account:

1. **Configure AWS credentials**:
   ```bash
   aws configure                    # Set up default profile
   # OR
   aws configure --profile myprof   # Set up named profile
   ```

2. **Run LazyCloud**:
   ```bash
   make run                        # Use default profile
   # OR  
   AWS_PROFILE=myprof make run    # Use specific profile
   ```

## Features

### Current (MVP)
- ✅ **Lambda Functions**: List, view details, environment variables
- ✅ **LocalStack Support**: Full development environment
- ✅ **Keyboard Navigation**: vim-like shortcuts
- ✅ **Real-time Updates**: Refresh with 'r'

### Coming Soon
- 🔄 **ECS Services**: Clusters, services, tasks, logs
- 🔄 **S3 Buckets**: Browse, view objects, basic operations
- 🔄 **EKS Clusters**: Status, nodes, basic workload info
- 🔄 **Log Viewer**: Integrated CloudWatch logs
- 🔄 **Multi-Region**: Easy region switching

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `q` | Quit application |
| `r` | Refresh current view |
| `?` | Show help |
| `j/k` or `↑/↓` | Navigate lists |
| `Enter` | Select item |

## Development

### Prerequisites
- Go 1.23+
- Docker (for LocalStack)
- Make

### Setup Development Environment
```bash
# Clone and setup
git clone <repo>
cd lazycloud

# Complete development setup (recommended)
make dev-setup

# Or manual setup
make dev-start     # Start LocalStack
make dev-run       # Run LazyCloud

# Other useful commands
make dev-status    # Check what test resources are available
make dev-reset     # Reset to clean state
make help          # See all commands
```

### Project Structure
```
lazycloud/
├── cmd/lazycloud/          # Main application
├── internal/
│   ├── aws/               # AWS service clients
│   │   ├── client.go      # Client manager
│   │   └── lambda/        # Lambda service
│   └── ui/views/          # TUI views
├── localstack/            # LocalStack setup
└── docker-compose.yml     # Development environment
```

### Testing with LocalStack

The project includes a complete LocalStack setup with test Lambda functions:

- **test-function**: Simple hello world
- **data-processor**: Processes data arrays
- **error-function**: Demonstrates error handling

### Development Commands

```bash
# Development workflow
make dev-setup     # Complete setup with test resources
make dev-run       # Run with LocalStack
make dev-status    # Check available test resources
make dev-stop      # Stop LocalStack
make dev-reset     # Reset to clean state

# Standard commands  
make build         # Build binary
make run           # Run with system AWS credentials
make test          # Run tests
make lint          # Run linter
make clean         # Clean build artifacts
make help          # Show all available commands
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `LAZYCLOUD_LOCAL` | Use LocalStack | `false` |
| `LOCALSTACK_ENDPOINT` | LocalStack URL | `http://localhost:4566` |
| `AWS_DEFAULT_REGION` | AWS region | `us-east-1` |

### AWS Authentication

LazyCloud uses the standard AWS credential chain:
1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. AWS credentials file (`~/.aws/credentials`)
3. IAM roles (for EC2/ECS)

## Architecture

LazyCloud follows a clean architecture with:
- **TUI Layer**: tview-based interface
- **Service Layer**: AWS service abstractions
- **Client Layer**: AWS SDK integration

Key design principles:
- **Demo-driven development**: Always maintain working demos
- **LocalStack first**: Develop against local AWS services
- **Security conscious**: Mask sensitive environment variables
- **Performance focused**: Efficient AWS API usage with caching

## Contributing

1. Check the [PROJECT_PLAN.md](PROJECT_PLAN.md) for current roadmap
2. See [DESIGN.md](DESIGN.md) for technical architecture
3. Review [CLAUDE.md](CLAUDE.md) for AI assistant context

## License

MIT License - see LICENSE file for details