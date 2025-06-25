#!/bin/bash

echo "🚀 Setting up LazyCloud test environment in LocalStack..."

# Create some test Lambda functions
echo "📦 Creating test Lambda functions..."

# Create a simple test function
cat > /tmp/test_function.py << 'EOF'
import json
import os
from datetime import datetime

def lambda_handler(event, context):
    print(f'Event: {json.dumps(event, indent=2)}')
    
    return {
        'statusCode': 200,
        'body': json.dumps({
            'message': 'Hello from LazyCloud test function!',
            'timestamp': datetime.utcnow().isoformat() + 'Z',
            'environment': os.environ.get('PYTHON_ENV', 'test'),
            'function_name': context.function_name if context else 'unknown'
        })
    }
EOF

# Create a more complex function
cat > /tmp/data_processor.py << 'EOF'
import json
import os
from datetime import datetime

def lambda_handler(event, context):
    print('Processing data...')
    
    # Simulate some processing
    data = event.get('data', [])
    processed = []
    
    for item in data:
        processed_item = {
            **item,
            'processed': True,
            'timestamp': datetime.utcnow().isoformat() + 'Z'
        }
        processed.append(processed_item)
    
    print(f'Processed {len(processed)} items')
    
    return {
        'statusCode': 200,
        'body': json.dumps({
            'message': 'Data processing complete',
            'processedCount': len(processed),
            'items': processed
        })
    }
EOF

# Create an error-prone function for testing error handling
cat > /tmp/error_function.py << 'EOF'
import json
from datetime import datetime

def lambda_handler(event, context):
    print('This function will demonstrate error handling')
    
    if event.get('shouldError'):
        raise Exception('Intentional error for testing error handling')
    
    return {
        'statusCode': 200,
        'body': json.dumps({
            'message': 'Function completed successfully',
            'timestamp': datetime.utcnow().isoformat() + 'Z'
        })
    }
EOF

# Zip the functions
cd /tmp
zip test-function.zip test_function.py
zip data-processor.zip data_processor.py
zip error-function.zip error_function.py

# Create Lambda functions
echo "Creating Lambda function: test-function"
awslocal lambda create-function \
    --function-name test-function \
    --runtime python3.9 \
    --role arn:aws:iam::123456789012:role/lambda-role \
    --handler test_function.lambda_handler \
    --zip-file fileb://test-function.zip \
    --description "Simple test function for LazyCloud development" \
    --timeout 30 \
    --memory-size 128 \
    --environment Variables='{PYTHON_ENV=development,APP_NAME=LazyCloud}'

echo "Creating Lambda function: data-processor"
awslocal lambda create-function \
    --function-name data-processor \
    --runtime python3.9 \
    --role arn:aws:iam::123456789012:role/lambda-role \
    --handler data_processor.lambda_handler \
    --zip-file fileb://data-processor.zip \
    --description "Data processing function for testing complex scenarios" \
    --timeout 60 \
    --memory-size 256 \
    --environment Variables='{PYTHON_ENV=development,BATCH_SIZE=100}'

echo "Creating Lambda function: error-function"
awslocal lambda create-function \
    --function-name error-function \
    --runtime python3.9 \
    --role arn:aws:iam::123456789012:role/lambda-role \
    --handler error_function.lambda_handler \
    --zip-file fileb://error-function.zip \
    --description "Function for testing error handling and logging" \
    --timeout 15 \
    --memory-size 64

# Create some S3 buckets
echo "🪣 Creating test S3 buckets..."
awslocal s3 mb s3://lazycloud-test-bucket
awslocal s3 mb s3://lazycloud-logs-bucket
awslocal s3 mb s3://lazycloud-data-bucket

# Add some test objects
echo "Adding test objects to S3..."
echo "Hello from LazyCloud!" > /tmp/test-file.txt
echo '{"message": "test data", "version": "1.0"}' > /tmp/test-data.json

awslocal s3 cp /tmp/test-file.txt s3://lazycloud-test-bucket/
awslocal s3 cp /tmp/test-data.json s3://lazycloud-test-bucket/data/
awslocal s3 cp /tmp/test_function.py s3://lazycloud-test-bucket/functions/

# Create some ECS resources (basic setup)
echo "🐳 Creating test ECS resources..."
awslocal ecs create-cluster --cluster-name lazycloud-test-cluster

# Create a simple task definition
cat > /tmp/task-definition.json << 'EOF'
{
    "family": "lazycloud-test-task",
    "networkMode": "bridge",
    "containerDefinitions": [
        {
            "name": "test-container",
            "image": "nginx:latest",
            "memory": 128,
            "portMappings": [
                {
                    "containerPort": 80,
                    "hostPort": 8080
                }
            ],
            "environment": [
                {
                    "name": "ENV",
                    "value": "test"
                }
            ]
        }
    ]
}
EOF

awslocal ecs register-task-definition --cli-input-json file:///tmp/task-definition.json

echo "✅ LocalStack setup complete!"
echo "🔗 Services available at: http://localhost:4566"
echo "📊 Health check: curl http://localhost:4566/_localstack/health"

# List created resources
echo "📋 Created resources:"
echo "Lambda functions:"
awslocal lambda list-functions --query 'Functions[].FunctionName' --output table

echo "S3 buckets:"
awslocal s3 ls

echo "ECS clusters:"
awslocal ecs list-clusters --query 'clusterArns' --output table