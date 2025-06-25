#!/bin/bash

echo "🚀 Setting up LazyCloud development environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start LocalStack
echo "🐳 Starting LocalStack..."
docker-compose up -d

# Wait for LocalStack to be ready
echo "⏳ Waiting for LocalStack to be ready..."
timeout=60
while [ $timeout -gt 0 ]; do
    if curl -f -s http://localhost:4566/_localstack/health > /dev/null 2>&1; then
        echo "✅ LocalStack is ready!"
        break
    fi
    echo "Waiting... ($timeout seconds remaining)"
    sleep 2
    timeout=$((timeout-2))
done

if [ $timeout -le 0 ]; then
    echo "❌ LocalStack failed to start within 60 seconds"
    echo "Check logs with: docker-compose logs localstack"
    exit 1
fi

# Show status
echo "📊 LocalStack status:"
curl -s http://localhost:4566/_localstack/health | jq .

echo ""
echo "🎉 Development environment ready!"
echo "🔗 LocalStack UI: http://localhost:4566"
echo "📚 AWS CLI commands use: --endpoint-url=http://localhost:4566"
echo ""
echo "Next steps:"
echo "  make run    # Start LazyCloud"
echo "  make test   # Run tests"
echo ""
echo "To stop LocalStack: docker-compose down"