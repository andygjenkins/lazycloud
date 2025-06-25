# LazyCloud - Technical Design Document

## Architecture Overview

LazyCloud follows a layered architecture pattern with clear separation between the TUI presentation layer, application logic, and AWS service integration.

```
┌─────────────────────────────────────────────────────┐
│                 TUI Layer                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │   Lambda    │ │     ECS     │ │     S3      │   │
│  │    View     │ │    View     │ │    View     │   │
│  └─────────────┘ └─────────────┘ └─────────────┘   │
│         │               │               │           │
└─────────┼───────────────┼───────────────┼───────────┘
          │               │               │
┌─────────┼───────────────┼───────────────┼───────────┐
│                Application Layer                    │
│  ┌─────────────────────────────────────────────────┐ │
│  │          Application State Manager              │ │
│  └─────────────────────────────────────────────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │   Lambda    │ │     ECS     │ │     S3      │   │
│  │   Service   │ │   Service   │ │   Service   │   │
│  └─────────────┘ └─────────────┘ └─────────────┘   │
└─────────┼───────────────┼───────────────┼───────────┘
          │               │               │
┌─────────┼───────────────┼───────────────┼───────────┐
│                AWS SDK Layer                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │AWS Lambda   │ │   AWS ECS   │ │   AWS S3    │   │
│  │   Client    │ │   Client    │ │   Client    │   │
│  └─────────────┘ └─────────────┘ └─────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Core Components Design

### 1. Application State Manager

**Purpose**: Central state management for the entire application.

**Key Responsibilities**:
- Manage current AWS region and profile
- Cache AWS resource data with TTL
- Handle global application state (current view, navigation history)
- Coordinate between different service views
- Manage real-time updates and refresh cycles

**Interface Design**:
```go
type AppState struct {
    // Authentication & Configuration
    CurrentProfile string
    CurrentRegion  string
    AWSConfig      aws.Config
    
    // Navigation State  
    CurrentView    ViewType
    ViewHistory    []ViewType
    
    // Resource Cache
    ResourceCache  *Cache
    
    // UI State
    KeyBindings    map[string]func()
    StatusMessage  string
    IsLoading      bool
}

type StateManager interface {
    // Configuration
    SetProfile(profile string) error
    SetRegion(region string) error
    GetAWSConfig() aws.Config
    
    // Navigation
    PushView(view ViewType)
    PopView() ViewType
    GetCurrentView() ViewType
    
    // Cache Management
    CacheResource(key string, data interface{}, ttl time.Duration)
    GetCachedResource(key string) (interface{}, bool)
    InvalidateCache(pattern string)
    
    // Status Updates
    SetStatus(message string)
    SetLoading(loading bool)
}
```

### 2. AWS Service Abstractions

Each AWS service will have a dedicated service wrapper that provides a simplified, opinionated interface tailored for the TUI.

#### Lambda Service Design

**Purpose**: Abstraction layer for AWS Lambda operations.

**Key Responsibilities**:
- List Lambda functions with pagination
- Retrieve function configuration and metadata
- Fetch CloudWatch logs for functions
- Execute test invocations
- Handle error cases gracefully

**Interface Design**:
```go
type LambdaFunction struct {
    Name         string
    Runtime      string
    Handler      string
    Memory       int64
    Timeout      int64
    LastModified time.Time
    Status       FunctionStatus
    Environment  map[string]string // Sensitive values masked
}

type LambdaService interface {
    // Function Management
    ListFunctions(ctx context.Context, region string) ([]*LambdaFunction, error)
    GetFunction(ctx context.Context, name string) (*LambdaFunction, error)
    
    // Logs
    GetRecentLogs(ctx context.Context, functionName string, hours int) ([]*LogEntry, error)
    StreamLogs(ctx context.Context, functionName string) (<-chan *LogEntry, error)
    
    // Testing
    InvokeFunction(ctx context.Context, functionName string, payload []byte) (*InvocationResult, error)
    
    // Monitoring
    GetMetrics(ctx context.Context, functionName string, duration time.Duration) (*FunctionMetrics, error)
}
```

#### ECS Service Design

**Purpose**: Abstraction for Amazon ECS clusters, services, and tasks.

**Interface Design**:
```go
type ECSCluster struct {
    Name               string
    Status             string
    RunningTasksCount  int
    PendingTasksCount  int
    ActiveServicesCount int
}

type ECSService struct {
    Name           string
    ClusterName    string
    Status         string
    DesiredCount   int
    RunningCount   int
    PendingCount   int
    TaskDefinition string
    LaunchType     string
}

type ECSServiceInterface interface {
    ListClusters(ctx context.Context) ([]*ECSCluster, error)
    ListServices(ctx context.Context, clusterName string) ([]*ECSService, error)
    ListTasks(ctx context.Context, clusterName, serviceName string) ([]*ECSTask, error)
    GetServiceLogs(ctx context.Context, clusterName, serviceName string) ([]*LogEntry, error)
}
```

### 3. TUI Components

#### Base Component Architecture

All UI components follow a consistent interface pattern:

```go
type Component interface {
    // Lifecycle
    Init() error
    Update(msg tea.Msg) tea.Cmd
    View() string
    
    // Focus Management
    Focus()
    Blur()
    Focused() bool
    
    // Key Handling
    HandleKey(key tea.KeyMsg) tea.Cmd
}

type ListComponent struct {
    items    []ListItem
    selected int
    focused  bool
    title    string
    height   int
}

type DetailComponent struct {
    content  string
    focused  bool
    title    string
    scrollY  int
}
```

#### View Management

Each AWS service has a dedicated view that manages the interaction between list and detail components:

```go
type LambdaView struct {
    functionList   *ListComponent
    functionDetail *DetailComponent
    logViewer      *LogComponent
    
    currentFunction *LambdaFunction
    state          ViewState
    service        LambdaService
}

type ViewState int
const (
    StateList ViewState = iota
    StateDetail
    StateLogs
)
```

### 4. Configuration Management

**Purpose**: Handle application configuration from multiple sources.

**Configuration Sources** (in priority order):
1. Command line flags
2. Environment variables
3. Configuration file (`~/.lazycloud/config.yaml`)
4. Defaults

**Configuration Structure**:
```yaml
# ~/.lazycloud/config.yaml
aws:
  default_profile: "default"
  default_region: "us-east-1"
  regions:
    - "us-east-1"
    - "us-west-2"
    - "eu-west-1"

ui:
  refresh_interval: 30s
  log_lines_limit: 1000
  theme: "dark"
  
keybindings:
  quit: "q"
  refresh: "r"
  help: "?"
  switch_region: "R"

services:
  lambda:
    enabled: true
    log_retention_hours: 24
  ecs:
    enabled: true
  s3:
    enabled: true
    max_objects_per_bucket: 1000
```

**Configuration Interface**:
```go
type Config struct {
    AWS         AWSConfig         `yaml:"aws"`
    UI          UIConfig          `yaml:"ui"`
    KeyBindings KeyBindingsConfig `yaml:"keybindings"`
    Services    ServicesConfig    `yaml:"services"`
}

type ConfigManager interface {
    Load() (*Config, error)
    Save(config *Config) error
    GetAWSProfiles() ([]string, error)
    ValidateConfig(config *Config) error
}
```

## Data Flow Architecture

### 1. Application Startup Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Parse     │───▶│   Load      │───▶│ Initialize  │
│   CLI Args  │    │   Config    │    │ AWS Clients │
└─────────────┘    └─────────────┘    └─────────────┘
        │                  │                  │
        ▼                  ▼                  ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Validate    │───▶│ Initialize  │───▶│   Start     │
│ AWS Creds   │    │ TUI App     │    │ Main Loop   │
└─────────────┘    └─────────────┘    └─────────────┘
```

### 2. Resource Loading Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ User Action │───▶│ Check       │───▶│ Load from   │
│ (Navigate)  │    │ Cache       │    │ AWS API     │
└─────────────┘    └─────────────┘    └─────────────┘
        │                  │                  │
        │                  ▼                  ▼
        │          ┌─────────────┐    ┌─────────────┐
        │          │ Return      │    │ Update      │
        │◀─────────│ Cached Data │    │ Cache       │
        │          └─────────────┘    └─────────────┘
        │                                     │
        ▼                                     ▼
┌─────────────┐                      ┌─────────────┐
│ Update UI   │◀─────────────────────│ Return      │
│ Components  │                      │ Fresh Data  │
└─────────────┘                      └─────────────┘
```

### 3. Error Handling Flow

```go
type ErrorHandler interface {
    HandleAWSError(err error) ErrorAction
    HandleUIError(err error) ErrorAction
    ShowUserError(message string, details string)
}

type ErrorAction int
const (
    ErrorActionRetry ErrorAction = iota
    ErrorActionIgnore
    ErrorActionExit
    ErrorActionPromptUser
)
```

## Performance Considerations

### 1. Caching Strategy

**Resource Cache**:
- TTL-based cache for AWS resource lists (default: 5 minutes)
- Aggressive caching for static data (IAM roles, region lists)
- Cache invalidation on user-triggered refresh
- Memory limits to prevent excessive usage

**Implementation**:
```go
type CacheEntry struct {
    Data      interface{}
    Timestamp time.Time
    TTL       time.Duration
}

type Cache struct {
    entries map[string]*CacheEntry
    mutex   sync.RWMutex
    maxSize int
}
```

### 2. Async Operations

**Background Loading**:
- Load AWS resources asynchronously to avoid blocking UI
- Show loading indicators for slow operations
- Implement cancellation for long-running operations

**Goroutine Management**:
```go
type TaskManager struct {
    tasks    map[string]context.CancelFunc
    mutex    sync.RWMutex
    maxTasks int
}

func (tm *TaskManager) StartTask(id string, task func(ctx context.Context)) {
    ctx, cancel := context.WithCancel(context.Background())
    tm.tasks[id] = cancel
    go task(ctx)
}
```

### 3. Memory Management

**Resource Limits**:
- Limit cached resources to prevent memory leaks
- Implement LRU eviction for cache entries
- Monitor memory usage and warn on high consumption

## Security Design

### 1. Credential Handling

**Principles**:
- Never store AWS credentials in application code
- Use AWS credential chain (environment, files, IAM roles)
- Support AWS profiles and role assumption
- No logging of sensitive information
- Hybrid approach: LocalStack for development, AWS credential chain for production

**Development vs Production Credential Strategy**:

**LocalStack Development**:
- Uses fixed test credentials (AWS_ACCESS_KEY_ID=test, AWS_SECRET_ACCESS_KEY=test)
- Detected via environment variables (LAZYCLOUD_LOCAL=true)
- Endpoint automatically set to http://localhost:4566
- No real AWS credentials needed

**Production Usage**:
- Follows standard AWS credential chain precedence:
  1. Environment variables
  2. AWS credentials file (~/.aws/credentials)
  3. AWS config file (~/.aws/config)
  4. IAM roles (for EC2/ECS/Lambda execution)
- Respects AWS_PROFILE environment variable
- Users expected to use `aws configure` for setup

**Implementation**:
```go
type CredentialManager struct {
    config aws.Config
    profiles map[string]aws.Config
    isLocalStack bool
}

func NewCredentialManager() (*CredentialManager, error) {
    ctx := context.Background()
    
    // Check if we're using LocalStack
    isLocalStack := os.Getenv("LOCALSTACK_ENDPOINT") != "" || 
        os.Getenv("AWS_ENDPOINT_URL") != "" ||
        os.Getenv("LAZYCLOUD_LOCAL") == "true"
    
    var cfg aws.Config
    var err error
    
    if isLocalStack {
        // Configure for LocalStack with test credentials
        cfg, err = config.LoadDefaultConfig(ctx,
            config.WithRegion("us-east-1"),
            config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
                func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                    endpoint := os.Getenv("LOCALSTACK_ENDPOINT")
                    if endpoint == "" {
                        endpoint = "http://localhost:4566"
                    }
                    return aws.Endpoint{
                        URL:           endpoint,
                        SigningRegion: region,
                    }, nil
                })),
        )
    } else {
        // Use standard AWS credential chain
        cfg, err = config.LoadDefaultConfig(ctx,
            config.WithRegion("us-east-1"), // default region
        )
    }
    
    if err != nil {
        return nil, err
    }
    
    return &CredentialManager{
        config: cfg,
        isLocalStack: isLocalStack,
    }, nil
}

func (cm *CredentialManager) LoadCredentials(profile string) (aws.Config, error) {
    if cm.isLocalStack {
        // Return LocalStack configuration
        return cm.config, nil
    }
    
    // Load specific AWS profile
    ctx := context.Background()
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithSharedConfigProfile(profile),
        config.WithRegion("us-east-1"), // default
    )
    return cfg, err
}
```

**User Experience Design**:

**Development Setup**:
```bash
# Simple LocalStack development
make dev-setup    # Auto-configures everything
make dev-run      # Uses test credentials automatically
```

**Production Setup**:
```bash
# Standard AWS approach
aws configure                 # Set up default credentials
make run                     # Uses system AWS credentials

# Or with specific profile
aws configure --profile prod  # Set up production profile
AWS_PROFILE=prod make run    # Use production profile
```

### 2. Data Sanitization

**Sensitive Data Handling**:
- Mask environment variables containing secrets
- Sanitize log outputs to remove credentials
- Hash or truncate sensitive identifiers in cache keys

## Testing Strategy

### 1. Unit Testing

**AWS Service Mocking**:
```go
type MockLambdaService struct {
    functions []*LambdaFunction
    logs      map[string][]*LogEntry
}

func (m *MockLambdaService) ListFunctions(ctx context.Context, region string) ([]*LambdaFunction, error) {
    return m.functions, nil
}
```

### 2. Integration Testing

**LocalStack Integration**:
- Use LocalStack for integration tests
- Test AWS service interactions without real AWS costs
- Validate error handling with simulated AWS errors

### 3. UI Testing

**TUI Testing Framework**:
```go
type UITestCase struct {
    Name     string
    Input    []tea.KeyMsg
    Expected string
    Setup    func(*testing.T) Component
}

func TestComponentBehavior(t *testing.T, testCase UITestCase) {
    component := testCase.Setup(t)
    for _, key := range testCase.Input {
        component.HandleKey(key)
    }
    assert.Equal(t, testCase.Expected, component.View())
}
```

## Deployment and Distribution

### 1. Build Process

**Cross-Platform Builds**:
```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/lazycloud-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o bin/lazycloud-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o bin/lazycloud-windows-amd64.exe
```

### 2. Distribution Channels

**Package Managers**:
- Homebrew (macOS/Linux)
- Scoop (Windows)
- APT/YUM repositories
- Direct GitHub releases

**Release Automation**:
- GitHub Actions for automated builds
- Semantic versioning
- Automated testing before release
- Binary signing and verification

This design provides a solid foundation for building LazyCloud with proper separation of concerns, testability, and maintainability while ensuring good performance and security practices.