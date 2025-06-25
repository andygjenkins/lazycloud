# LazyCloud - AWS Terminal UI Project Plan & Implementation Guide

## Project Overview

LazyCloud is a terminal-based user interface (TUI) application for managing and monitoring AWS resources, inspired by tools like lazydocker, lazygit, and k9s. The goal is to provide developers and DevOps engineers with a fast, keyboard-driven interface to view AWS resources, check statuses, and access logs without leaving the terminal.

## Vision Statement

"Make AWS resource management as intuitive and efficient as Git management with lazygit - providing instant visibility into cloud infrastructure through a beautiful, responsive terminal interface."

## Core Principles

1. **Keyboard-First**: All navigation and actions should be keyboard-driven
2. **Performance**: Fast loading and responsive interface, even with large numbers of resources
3. **Simplicity**: Clean, uncluttered interface focusing on essential information
4. **Extensibility**: Architecture that allows easy addition of new AWS services
5. **Security**: Secure credential handling and no sensitive data logging
6. **Demo-Driven Development**: Always maintain a working demo

## Early Demo & Testing Philosophy

**Core Principle**: Always have a working demo, no matter how basic.

### Demo-Driven Development Approach:
- **Day 2**: Basic TUI runs and displays "Hello LazyCloud"
- **Day 3**: Connects to AWS and shows real Lambda data
- **Every 2-3 days**: New working feature demonstrated
- **End of each week**: Fully functional subset of features

### Testing Strategy:
- **Manual Testing First**: Each demo checkpoint manually verified
- **Unit Tests**: Added incrementally as components stabilize
- **Integration Tests**: Added once basic AWS integration works
- **Continuous Testing**: Every commit should maintain working demo

## Target Users & Use Cases

- **Primary**: DevOps engineers and developers working with AWS daily
- **Secondary**: SREs, platform engineers, and cloud architects
- **Use Cases**: 
  - Quick status checks of deployments
  - Log viewing and troubleshooting
  - Resource discovery and exploration
  - Pre-deployment verification

## MVP Scope

### Core Services
1. **AWS Lambda** - List functions, view config, access logs, test invocations
2. **Amazon ECS** - List clusters/services/tasks, view definitions, access container logs
3. **Amazon EKS** - List clusters, view status, basic pod listing, cluster logs
4. **Amazon S3** - List buckets, browse contents, view metadata, basic operations

### Core Features
- Multi-region support with region switching
- Real-time status updates
- Integrated log viewer with search/filtering
- Credential management (AWS profiles, IAM roles)
- Configuration persistence
- Help system and keyboard shortcuts

## Technical Architecture

### Technology Stack
- **Language**: Go (chosen for AWS SDK quality, performance, single binary distribution)
- **TUI Framework**: tview for MVP (mature, widget-based), consider bubbletea for v2
- **AWS Integration**: AWS SDK for Go v2
- **Configuration**: YAML files + environment variables
- **Build**: Single binary with cross-platform support

### Project Structure
```
lazycloud/
├── cmd/lazycloud/          # Main application entry point
├── internal/
│   ├── app/                # Application core and state management
│   ├── ui/                 # TUI components and layouts
│   │   ├── components/     # Reusable UI components
│   │   └── views/          # Service-specific views
│   ├── aws/                # AWS service abstractions
│   │   ├── lambda/
│   │   ├── ecs/
│   │   ├── eks/
│   │   └── s3/
│   ├── config/             # Configuration management
│   └── utils/              # Shared utilities
├── pkg/types/              # Public types and interfaces
├── docs/                   # Documentation
└── scripts/                # Build and deployment scripts
```

### Key Bindings (Following lazygit conventions)
- `tab`/`shift+tab`: Navigate between panels
- `j/k` or `↑/↓`: Navigate lists
- `enter`: Select/drill down
- `esc`: Go back/cancel
- `r`: Refresh current view
- `R`: Change region
- `q`: Quit application
- `/`: Search within current view
- `?`: Help/shortcuts

## Development Workflow

### LocalStack Development Setup

LazyCloud includes a complete LocalStack development environment with pre-configured AWS resources for testing:

- **Lambda Functions**: 3 test functions (test-function, data-processor, error-function)
- **S3 Buckets**: 3 test buckets with sample objects  
- **ECS Resources**: Test cluster and task definitions
- **Automatic Setup**: Running lambdas and resources ready immediately

#### Development Commands

**LocalStack Development (Recommended)**
```bash
# Complete development setup (first time)
make dev-setup     # Starts LocalStack with test Lambda functions, S3 buckets, ECS resources

# Daily development workflow
make dev-start     # Start LocalStack
make dev-run       # Run LazyCloud with LocalStack credentials
make dev-stop      # Stop LocalStack
make dev-status    # Check what resources are available

# Reset development environment
make dev-reset     # Clean slate with fresh test data
```

**Standard Development Commands**
```bash
# Build and run
make build         # Build binary to bin/lazycloud
make run           # Run application (uses AWS credentials from system)

# Testing and quality
make test          # Run tests
make lint          # Run linter
make deps          # Install/update dependencies

# Clean up
make clean         # Remove build artifacts
```

**AWS Credential Management**
```bash
# For production AWS usage
aws configure                    # Set up AWS credentials
aws configure --profile dev      # Set up development profile
AWS_PROFILE=dev make run        # Use specific profile

# For LocalStack development
make dev-run                    # Uses test credentials automatically
```

## Implementation Status

### ✅ Development Environment (Completed)
- [x] LocalStack integration with Docker Compose
- [x] Makefile with development workflow commands  
- [x] Test Lambda functions and S3 buckets automatically created
- [x] Credential management (LocalStack vs AWS)
- [x] Environment configuration (.env.example)

## Implementation Plan

### Phase 1: Foundation (Week 1) - EARLY DEMO FOCUSED

#### Demo Milestones:
- **Day 2**: Basic TUI app runs and displays "Hello LazyCloud"
- **Day 3**: AWS credentials working, can list Lambda functions
- **Day 5**: Full Lambda list view with basic navigation
- **Day 7**: Lambda detail view showing function info

#### Tasks by Priority:

**Days 1-2: Minimal Viable Setup**
- [x] Initialize Go module: `go mod init lazycloud`
- [x] Create minimal directory structure (`cmd/`, `internal/`)
- [x] Set up LocalStack development environment with Docker
- [x] Create enhanced `Makefile` with development workflow
- [x] Set up credential management and environment configuration
- [ ] Quick TUI framework POC comparison (tview vs bubbletea)
- [ ] Create basic `cmd/lazycloud/main.go` with "Hello LazyCloud" TUI
- [ ] **DEMO CHECKPOINT**: App runs and shows basic TUI

**Days 3-4: AWS Integration Proof of Concept**
- [x] Add AWS SDK dependencies to go.mod
- [x] Create AWS client setup with LocalStack/production switching
- [ ] Implement simple Lambda function listing (just names)
- [ ] Display Lambda list in TUI (basic list view)
- [ ] **DEMO CHECKPOINT**: Shows real Lambda functions from LocalStack
  - Test data available: test-function, data-processor, error-function

**Days 4-5: Configuration Management**
- [ ] Create `internal/config/config.go` with basic struct
- [ ] Implement AWS credential chain handling
- [ ] Add region switching functionality
- [ ] Create command-line flags using cobra

**Days 5-7: Enhanced Lambda View**
- [ ] Add function details panel
- [ ] Implement basic log viewer
- [ ] Add keyboard navigation between panels
- [ ] **DEMO CHECKPOINT**: Complete Lambda viewer working

### Phase 2: Core Services (Weeks 2-3)
- [ ] Complete Lambda functionality (logs, details, invoke)
- [ ] ECS service integration
- [ ] S3 basic functionality
- [ ] Multi-region support
- [ ] Tabbed navigation between services

### Phase 3: Polish & EKS (Weeks 4-5)
- [ ] EKS integration
- [ ] Enhanced log viewing (search, filtering)
- [ ] Caching and performance optimizations
- [ ] Error handling and user feedback
- [ ] Help system

### Phase 4: MVP Release (Week 6)
- [ ] Testing and bug fixes
- [ ] Documentation and help system
- [ ] Build and release automation
- [ ] Basic packaging (Homebrew, releases)

## Success Metrics

### Technical Metrics
- Application startup time < 2 seconds
- API response handling < 500ms for most operations
- Memory usage < 50MB under normal operation
- Support for AWS accounts with 100+ resources per service

### User Experience Metrics
- New user can navigate basic functionality within 5 minutes
- Common tasks (checking Lambda logs) achievable in < 3 keystrokes
- Each demo milestone achieved on schedule

## Risk Assessment

### High Risk
- **AWS API Rate Limits**: Mitigate with intelligent caching and request batching
- **Complex EKS Integration**: May require kubectl dependencies or simplified view

### Medium Risk
- **Cross-platform Terminal Compatibility**: Test across different terminal emulators
- **Large Account Performance**: Implement pagination and lazy loading

### Low Risk
- **Credential Security**: Leverage existing AWS credential chain
- **Configuration Complexity**: Start simple, add features incrementally

## Future Enhancements (Post-MVP)

### Additional Services
- CloudWatch (metrics and alarms)
- RDS (database instances and logs)
- EC2 (instances and security groups)
- CloudFormation/CDK stacks

### Advanced Features
- Resource tagging and filtering
- Custom dashboards
- Integration with AWS Organizations
- Cost optimization insights
- Plugin system for custom services

## Feedback & Implementation Notes

### User Comments/Suggestions:
*Add your feedback and suggestions here as we implement*

### Demo Results Log:
*Track results of each demo checkpoint*
- [ ] Day 2 Demo: _______________
- [ ] Day 3 Demo: _______________
- [ ] Day 5 Demo: _______________
- [ ] Week 1 Final: _____________

### Implementation Decisions:
*Record key technical decisions made during development*

### Phase Adjustments:
*Document changes to the plan as we learn during implementation*

---

**Next Steps**: Ready to start Phase 1, Day 1 - Initialize Go project and create basic TUI structure.