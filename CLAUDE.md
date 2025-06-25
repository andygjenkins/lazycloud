# LazyCloud - AI Assistant Context

## Project Overview
LazyCloud is a terminal-based user interface (TUI) application for managing AWS resources, similar to lazydocker, lazygit, and k9s. It provides a keyboard-driven interface for viewing AWS services like Lambda, ECS, EKS, and S3 with integrated log viewing and status monitoring.

## Current Status
- **Phase**: Planning and Design
- **Language**: Go (decided)
- **TUI Framework**: TBD (tview vs bubbletea)
- **Target MVP Services**: Lambda, ECS, EKS, S3

## Technology Stack

### Core Technologies
- **Language**: Go (chosen for AWS SDK quality, performance, single binary distribution)
- **TUI Framework**: Likely tview for MVP (mature, widget-based, easier complex layouts)
- **AWS Integration**: AWS SDK for Go v2
- **Configuration**: YAML/TOML files + environment variables
- **Build**: Single binary with cross-platform support

### Key Dependencies (Planned)
```go
// AWS SDK
github.com/aws/aws-sdk-go-v2
github.com/aws/aws-sdk-go-v2/service/lambda
github.com/aws/aws-sdk-go-v2/service/ecs  
github.com/aws/aws-sdk-go-v2/service/eks
github.com/aws/aws-sdk-go-v2/service/s3

// TUI Framework (one of these)
github.com/rivo/tview                    // Likely choice for MVP
github.com/charmbracelet/bubbletea      // Future consideration

// Configuration
github.com/spf13/viper
github.com/spf13/cobra                   // CLI framework
```

## Project Structure
```
lazycloud/
├── cmd/lazycloud/          # Main application entry point
├── internal/
│   ├── app/                # Application core and state management  
│   ├── ui/                 # TUI components and layouts
│   │   ├── components/     # Reusable UI components
│   │   └── views/          # Service-specific views (lambda, ecs, etc.)
│   ├── aws/                # AWS service abstractions
│   │   ├── lambda/         # Lambda service client wrapper
│   │   ├── ecs/            # ECS service client wrapper
│   │   ├── eks/            # EKS service client wrapper
│   │   └── s3/             # S3 service client wrapper
│   ├── config/             # Configuration management
│   └── utils/              # Shared utilities
├── pkg/types/              # Public types and interfaces
├── docs/                   # Documentation
└── scripts/                # Build and deployment scripts
```

## Development Workflow

### Common Commands
```bash
# Initialize Go module (when ready)
go mod init lazycloud

# Run the application
go run cmd/lazycloud/main.go

# Build binary
go build -o bin/lazycloud cmd/lazycloud/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint (install golangci-lint first)
golangci-lint run

# Tidy dependencies
go mod tidy
```

### Key Development Files
- `PROJECT_PLAN.md` - Comprehensive project planning document
- `DESIGN.md` - Technical architecture and design decisions (TBD)
- `cmd/lazycloud/main.go` - Application entry point
- `internal/app/app.go` - Core application logic and state
- `internal/config/config.go` - Configuration management
- `go.mod` - Go module definition and dependencies

## Code Style and Conventions

### Go Standards
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Prefer composition over inheritance
- Handle errors explicitly, don't ignore them
- Use interfaces for testability and modularity

### Project-Specific Conventions
- AWS service wrappers in `internal/aws/` should provide simplified, opinionated interfaces
- UI components should be stateless where possible, with state managed in `internal/app/`
- Configuration should be environment-aware (dev/staging/prod)
- All AWS API calls should include proper error handling and retries
- Sensitive information (credentials, keys) should never be logged

### Package Organization
- `internal/` for application-specific code not intended for external use
- `pkg/` for code that could be reused in other projects
- Service-specific code grouped by AWS service
- UI code separated from business logic

## AWS Integration Notes

### Authentication Strategy
- Use AWS credential chain (environment variables, ~/.aws/credentials, IAM roles)
- Support multiple AWS profiles
- Allow region selection/switching within the application
- No hardcoded credentials ever

### API Usage Patterns
- Use AWS SDK v2 context-aware APIs
- Implement proper retries and exponential backoff
- Cache frequently accessed data (with TTL)
- Handle rate limiting gracefully
- Support pagination for large result sets

### Supported Regions
- Start with common regions (us-east-1, us-west-2, eu-west-1)
- Allow dynamic region switching
- Handle region-specific service availability

## Testing Strategy

### Unit Tests
- Test AWS service wrappers with mocked AWS clients
- Test UI components with simulated input
- Test configuration parsing and validation
- Coverage target: >80% for core business logic

### Integration Tests
- Test against LocalStack for AWS services (where possible)
- Test with real AWS resources in dedicated test account
- Test different AWS credential scenarios

### Manual Testing
- Test in different terminal environments
- Test with different AWS account sizes
- Test keyboard navigation and shortcuts
- Test error conditions and edge cases

## UI/UX Guidelines

### Keyboard Navigation
- Follow vim-like conventions where applicable
- `j/k` for list navigation
- `tab/shift+tab` for panel switching
- `enter` for selection/drill-down
- `esc` for back/cancel
- `q` for quit
- `?` for help

### Visual Design
- Dark theme by default
- Consistent color coding (green=healthy, red=error, yellow=warning)
- Clear visual hierarchy
- Minimal information density - show what's essential
- Status indicators for loading/error states

### Error Handling
- User-friendly error messages
- Graceful degradation when services are unavailable
- Clear indication of network/authentication issues
- Retry mechanisms with user feedback

## Security Considerations
- Never log AWS credentials or sensitive data
- Use secure credential storage mechanisms
- Validate all user inputs
- Follow AWS security best practices
- Regular dependency security updates

## Common Tasks for AI Assistants

### When adding a new AWS service:
1. Create service wrapper in `internal/aws/servicename/`
2. Add configuration options in `internal/config/`
3. Create UI views in `internal/ui/views/servicename/`
4. Add navigation in main app
5. Update help documentation
6. Add tests for new functionality

### When debugging issues:
1. Check AWS credential configuration
2. Verify service availability in target region
3. Check network connectivity and permissions
4. Review application logs for error details
5. Test with minimal AWS resources first

### When optimizing performance:
1. Profile memory usage with large accounts
2. Check AWS API call patterns for efficiency
3. Implement proper caching strategies
4. Optimize UI rendering for large datasets
5. Consider lazy loading for expensive operations

## Future Considerations
- Plugin system for custom AWS services
- Integration with AWS Organizations
- Export capabilities (JSON, CSV)
- Web companion interface
- IDE integrations

## Notes for AI Assistants
- Always check `PROJECT_PLAN.md` for current project scope and decisions
- Follow the established project structure strictly
- Test any AWS-related code thoroughly
- Consider terminal compatibility across different platforms
- Prioritize keyboard navigation and accessibility
- Keep security and credential handling as top priority