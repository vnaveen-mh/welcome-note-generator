#!/bin/bash

# Welcome Note Generator - Web UI Startup Script

echo "🚀 Welcome Note Generator - Starting Web Server"
echo "================================================"

# Check if GEMINI_API_KEY is set
if [ -z "$GEMINI_API_KEY" ]; then
    echo "❌ Error: GEMINI_API_KEY environment variable is not set"
    echo ""
    echo "Please set your Google API key:"
    echo "  export GEMINI_API_KEY=your_api_key_here"
    echo ""
    echo "Or run this script with:"
    echo "  GEMINI_API_KEY=your_key ./run-web.sh"
    exit 1
fi

echo "✓ API key configured"

# Generate TEMPL files
echo ""
echo "📝 Generating TEMPL templates..."
if command -v templ &> /dev/null; then
    templ generate
    echo "✓ TEMPL files generated"
else
    echo "⚠️  Warning: templ command not found. Using pre-generated files."
    echo "   Install templ with: go install github.com/a-h/templ/cmd/templ@latest"
fi

# Build and run
echo ""
echo "🔨 Building application..."
go build -o welcome-note-web cmd/web/main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful"
    echo ""
    echo "🌐 Starting web server on http://localhost:8080"
    echo "================================================"
    echo ""
    ./welcome-note-web
else
    echo "❌ Build failed"
    exit 1
fi
