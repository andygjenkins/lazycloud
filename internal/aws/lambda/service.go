package lambda

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type Service struct {
	client *lambda.Client
}

type Function struct {
	Name         string
	Runtime      string
	Handler      string
	Description  string
	Memory       int32
	Timeout      int32
	LastModified time.Time
	Status       string
	Environment  map[string]string
}

func NewService(client *lambda.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ListFunctions(ctx context.Context) ([]*Function, error) {
	var functions []*Function
	
	input := &lambda.ListFunctionsInput{}
	
	// Use paginator to handle multiple pages
	paginator := lambda.NewListFunctionsPaginator(s.client, input)
	
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		
		for _, fn := range page.Functions {
			function := &Function{
				Name:        *fn.FunctionName,
				Runtime:     string(fn.Runtime),
				Handler:     *fn.Handler,
				Memory:      *fn.MemorySize,
				Timeout:     *fn.Timeout,
				Status:      string(fn.State),
				Environment: make(map[string]string),
			}
			
			if fn.Description != nil {
				function.Description = *fn.Description
			}
			
			// Parse last modified time
			if fn.LastModified != nil {
				if t, err := time.Parse(time.RFC3339, *fn.LastModified); err == nil {
					function.LastModified = t
				}
			}
			
			// Get environment variables (mask sensitive ones)
			if fn.Environment != nil && fn.Environment.Variables != nil {
				for k, v := range fn.Environment.Variables {
					if isSensitiveEnvVar(k) {
						function.Environment[k] = "***masked***"
					} else {
						function.Environment[k] = v
					}
				}
			}
			
			functions = append(functions, function)
		}
	}
	
	return functions, nil
}

func (s *Service) GetFunction(ctx context.Context, name string) (*Function, error) {
	input := &lambda.GetFunctionInput{
		FunctionName: &name,
	}
	
	result, err := s.client.GetFunction(ctx, input)
	if err != nil {
		return nil, err
	}
	
	fn := result.Configuration
	function := &Function{
		Name:        *fn.FunctionName,
		Runtime:     string(fn.Runtime),
		Handler:     *fn.Handler,
		Memory:      *fn.MemorySize,
		Timeout:     *fn.Timeout,
		Status:      string(fn.State),
		Environment: make(map[string]string),
	}
	
	if fn.Description != nil {
		function.Description = *fn.Description
	}
	
	// Parse last modified time
	if fn.LastModified != nil {
		if t, err := time.Parse(time.RFC3339, *fn.LastModified); err == nil {
			function.LastModified = t
		}
	}
	
	// Get environment variables (mask sensitive ones)
	if fn.Environment != nil && fn.Environment.Variables != nil {
		for k, v := range fn.Environment.Variables {
			if isSensitiveEnvVar(k) {
				function.Environment[k] = "***masked***"
			} else {
				function.Environment[k] = v
			}
		}
	}
	
	return function, nil
}

func (s *Service) InvokeFunction(ctx context.Context, name string, payload []byte) (*InvocationResult, error) {
	input := &lambda.InvokeInput{
		FunctionName: &name,
		Payload:      payload,
	}
	
	result, err := s.client.Invoke(ctx, input)
	if err != nil {
		return nil, err
	}
	
	invocationResult := &InvocationResult{
		StatusCode: result.StatusCode,
		Payload:    result.Payload,
	}
	
	if result.FunctionError != nil {
		invocationResult.Error = *result.FunctionError
	}
	
	if result.LogResult != nil {
		invocationResult.LogResult = *result.LogResult
	}
	
	return invocationResult, nil
}

type InvocationResult struct {
	StatusCode int32
	Payload    []byte
	Error      string
	LogResult  string
}

// Helper function to determine if an environment variable is sensitive
func isSensitiveEnvVar(key string) bool {
	sensitiveKeys := []string{
		"PASSWORD", "PASSWD", "SECRET", "KEY", "TOKEN", "API_KEY",
		"AWS_SECRET_ACCESS_KEY", "DATABASE_PASSWORD", "DB_PASSWORD",
		"PRIVATE_KEY", "CERT", "CREDENTIAL",
	}
	
	keyUpper := key
	for _, sensitive := range sensitiveKeys {
		if keyUpper == sensitive || len(keyUpper) > len(sensitive) {
			// Simple substring check for common patterns
			for i := 0; i <= len(keyUpper)-len(sensitive); i++ {
				if keyUpper[i:i+len(sensitive)] == sensitive {
					return true
				}
			}
		}
	}
	
	return false
}