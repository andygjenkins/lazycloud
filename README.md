# LazyCloud 🚀

AWS Terminal UI for the Modern Developer - inspired by lazydocker, lazygit, and k9s.

## Quick Start

### With LocalStack (Development)

1. **Start LocalStack**:
   ```bash
   docker-compose up -d
   ```

2. **Run LazyCloud**:
   ```bash
   LAZYCLOUD_LOCAL=true make run
   ```

### With Real AWS

1. **Configure AWS credentials** (one of these methods):
   ```bash
   # Using AWS CLI
   aws configure
   
   # Using environment variables
   export AWS_ACCESS_KEY_ID=your_key
   export AWS_SECRET_ACCESS_KEY=your_secret
   export AWS_DEFAULT_REGION=us-east-1
   ```

2. **Run LazyCloud**:
   ```bash
   make run
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

# Start LocalStack with test data
docker-compose up -d

# Build and run
make build
LAZYCLOUD_LOCAL=true ./bin/lazycloud
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

### Build Commands

```bash
make build          # Build binary
make run           # Run with go run
make test          # Run tests
make clean         # Clean build artifacts
make build-all     # Build for all platforms
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