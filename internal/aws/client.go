package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ClientManager struct {
	config aws.Config
	region string
	profile string
	
	// Service clients
	lambdaClient *lambda.Client
	s3Client     *s3.Client
	ecsClient    *ecs.Client
}

func NewClientManager() (*ClientManager, error) {
	ctx := context.Background()
	
	// Check if we're using LocalStack
	isLocalStack := os.Getenv("LOCALSTACK_ENDPOINT") != "" || 
		os.Getenv("AWS_ENDPOINT_URL") != "" ||
		os.Getenv("LAZYCLOUD_LOCAL") == "true"
	
	var cfg aws.Config
	var err error
	
	if isLocalStack {
		// Configure for LocalStack
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
		// Configure for real AWS
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-east-1"), // default region
		)
	}
	
	if err != nil {
		return nil, err
	}
	
	cm := &ClientManager{
		config: cfg,
		region: cfg.Region,
	}
	
	// Initialize service clients
	cm.lambdaClient = lambda.NewFromConfig(cfg)
	cm.s3Client = s3.NewFromConfig(cfg)
	cm.ecsClient = ecs.NewFromConfig(cfg)
	
	return cm, nil
}

func (cm *ClientManager) GetLambdaClient() *lambda.Client {
	return cm.lambdaClient
}

func (cm *ClientManager) GetS3Client() *s3.Client {
	return cm.s3Client
}

func (cm *ClientManager) GetECSClient() *ecs.Client {
	return cm.ecsClient
}

func (cm *ClientManager) GetRegion() string {
	return cm.region
}

func (cm *ClientManager) SetRegion(region string) error {
	// Update config with new region
	cfg := cm.config.Copy()
	cfg.Region = region
	
	// Recreate clients with new region
	cm.lambdaClient = lambda.NewFromConfig(cfg)
	cm.s3Client = s3.NewFromConfig(cfg)
	cm.ecsClient = ecs.NewFromConfig(cfg)
	
	cm.region = region
	cm.config = cfg
	
	return nil
}

func (cm *ClientManager) TestConnection(ctx context.Context) error {
	// Test connection by trying to list Lambda functions
	_, err := cm.lambdaClient.ListFunctions(ctx, &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(1),
	})
	return err
}