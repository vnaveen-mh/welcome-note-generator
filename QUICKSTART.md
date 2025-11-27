# Welcome Note Generator - Quick Start Guide

## Overview

This project demonstrates building production-ready AI applications with Go and Genkit, featuring a stunning web UI powered by Datastar and TEMPL.

## What You Get

- **5 Progressive Flow Versions**: From simple to advanced AI workflows
- **Beautiful Web UI**: Single-page application with tabs for each flow
- **Real-time Updates**: Datastar-powered hypermedia interactions
- **Content Safety**: Built-in moderation pipeline
- **Natural Language**: Smart flow interprets free-form requests
- **Production Ready**: Rate limiting, CSRF protection, Docker deployment

## Prerequisites

1. **Go 1.25+** installed
2. **Gemini API Key** for Google Gemini
3. **Docker & Docker Compose** (optional, for containerized deployment)
4. **templ CLI** (optional, for development)

## Quick Start

Choose your preferred deployment method:

### Option 1: Docker Deployment (Recommended)

**Fastest way to get started with production features enabled:**

1. **Get Your API Key**

   Sign up at [Google AI Studio](https://makersuite.google.com/app/apikey) and get your API key.

2. **Set Up Environment**

   ```bash
   cp .env.example .env
   # Edit .env and add your GEMINI_API_KEY
   ```

3. **Start with Docker Compose**

   ```bash
   docker-compose up -d
   ```

4. **Open Your Browser**

   Navigate to: **http://localhost:8080**

See [DOCKER.md](DOCKER.md) for comprehensive Docker deployment documentation.

---

### Option 2: Local Development

**For development and customization:**

1. **Get Your API Key**

   Sign up at [Google AI Studio](https://makersuite.google.com/app/apikey) and get your API key.

2. **Set Environment Variable**

   ```bash
   export GEMINI_API_KEY=your_api_key_here
   ```

3. **Install Dependencies**

   ```bash
   go mod download
   ```

4. **Run the Web Server**

   **Option A: Using the startup script**

   ```bash
   ./run-web.sh
   ```

   **Option B: Directly with go run**

   ```bash
   # Generate TEMPL files first
   go run github.com/a-h/templ/cmd/templ@latest generate

   # Run the server
   go run cmd/web/main.go
   ```

5. **Open Your Browser**

   Navigate to: **http://localhost:8080**

---

### Option 3: Genkit Developer UI

**For flow debugging and tracing:**

```bash
# Run Genkit developer UI
go run cmd/genkit/main.go

# Access Genkit UI
# http://localhost:4000 (or check console output)
```

## Environment Variables

Configure the application using environment variables:

| Variable                         | Description                     | Default        | Required |
| -------------------------------- | ------------------------------- | -------------- | -------- |
| `GEMINI_API_KEY`                 | Google Gemini API key           | -              | Yes      |
| `PORT`                           | Server port                     | `8080`         | No       |
| `CSRF_KEY`                       | 32-byte CSRF key (hex)          | Auto-generated | No       |
| `CSRF_TRUSTED_ORIGINS`           | Comma-separated trusted origins | Empty          | No       |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Rate limit per IP               | `30`           | No       |
| `RATE_LIMIT_BURST_SIZE`          | Rate limit burst size           | `5`            | No       |
| `RATE_LIMIT_CLEANUP_INTERVAL`    | Cleanup interval                | `5m`           | No       |
| `RATE_LIMIT_LIMITER_TTL`         | Limiter TTL                     | `15m`          | No       |

**Example `.env` file:**

```bash
GEMINI_API_KEY=your_gemini_api_key_here
PORT=8080
CSRF_KEY=<yout_csrf_key>
CSRF_TRUSTED_ORIGINS=http://localhost:8080,https://yourdomain.com
RATE_LIMIT_REQUESTS_PER_MINUTE=30
RATE_LIMIT_BURST_SIZE=5
```

## Using the UI

### Tab Navigation

Click on any tab to switch between flow versions:

1. **V1: Simple** - Basic string input/output
2. **V2: Structured** - Rich customization options
3. **V3: Metadata** - Includes generation details
4. **Safe Flow** - With content moderation
5. **Smart Flow** - Natural language input

### Example Usage

#### V1: Simple Flow

1. Click "V1: Simple" tab
2. Enter: `birthday party`
3. Click "Generate Welcome Note"
4. See your personalized welcome note

#### V2: Structured Input

1. Click "V2: Structured" tab
2. Fill in:
   - Occasion: `hotel check-in`
   - Language: `Spanish`
   - Length: `short`
   - Tone: `warm`
3. Click "Generate Customized Note"
4. Get a Spanish welcome note for hotel guests

#### V3: Structured Output

1. Click "V3: Metadata" tab
2. Enter occasion: `new employee onboarding`
3. Select your preferences
4. Click "Generate with Metadata"
5. See the note plus all generation parameters

#### Safe Flow: Content Moderation

1. Click "Safe Flow" tab
2. Try entering:
   - Normal request: `Christmas party`
   - Test moderation: Use tone `insulting` or `sarcastic`
3. Watch the moderation system sanitize inappropriate content
4. See moderation status and notes

#### Smart Flow: Natural Language

1. Click "Smart Flow" tab
2. Describe what you want in plain English:
   ```
   Write a warm and professional welcome message for our new
   software engineer joining next Monday. Make it friendly
   but not too casual.
   ```
3. The AI extracts parameters automatically
4. Generates and moderates the note
5. Returns complete result with metadata

## Project Structure

```
welcome-note-generator/
├── cmd/
│   ├── web/
│   │   └── main.go              # Web server entry point
│   └── genkit/
│       └── main.go              # Genkit developer UI
├── internal/
│   ├── flows/
│   │   ├── v1.go                # Simple flow
│   │   ├── v2.go                # Structured input flow
│   │   ├── v3.go                # Structured output flow
│   │   ├── safe_flow.go         # Moderation pipeline
│   │   ├── smart_flow.go        # NLP interpretation flow
│   │   ├── toxicity_filter.go   # Content moderation
│   │   ├── flow_store.go        # Flow registry
│   │   └── tones.go             # Tone definitions
│   └── types/
│       └── types.go             # Type definitions
├── web/
│   ├── handlers/                # HTTP request handlers
│   │   ├── v1_handler.go
│   │   ├── v2_handler.go
│   │   ├── v3_handler.go
│   │   ├── safe_handler.go
│   │   └── smart_handler.go
│   ├── middleware/              # HTTP middleware
│   │   ├── rate-limit.go        # Per-IP rate limiting
│   │   ├── csrf.go              # CSRF protection
│   │   └── logger.go            # Structured logging
│   ├── templates/               # Templ components
│   │   ├── layout.templ
│   │   ├── index.templ
│   │   ├── layout_templ.go      # Generated
│   │   └── index_templ.go       # Generated
│   ├── utils/
│   │   └── datastar_helpers.go  # SSE utilities
│   └── config/
│       └── config.go            # Configuration management
├── Dockerfile                   # Multi-stage Alpine build
├── docker-compose.yml           # Docker deployment
├── .env.example                 # Environment template
├── README.md                    # Project overview
├── DOCKER.md                    # Docker documentation
├── QUICKSTART.md                # This file
└── article_final.md             # Medium article
```

## Production Features

### Rate Limiting

The application includes per-IP rate limiting on all `/api/*` endpoints:

- **Default**: 30 requests per minute per IP
- **Burst**: 5 additional requests allowed
- **Configurable**: Via environment variables (no rebuild needed)

When rate limit is exceeded, you'll see a user-friendly error in the UI.

### CSRF Protection

All form submissions are protected with CSRF tokens:

- Auto-generated secure tokens
- SameSite cookie protection
- Configurable trusted origins
- Custom middleware handles trailing slash inconsistencies

### Security Headers

- Secure cookie settings
- HTTPS-ready configuration
- Origin validation

## API Endpoints

If you prefer using the API directly:

### Web UI Endpoints

All `/api/*` endpoints include rate limiting and CSRF protection:

- `GET /` - Main UI
- `POST /api/v1/generate` - V1 flow
- `POST /api/v2/generate` - V2 flow
- `POST /api/v3/generate` - V3 flow
- `POST /api/safe/generate` - Safe flow
- `POST /api/smart/generate` - Smart flow

## Testing Scenarios

### Birthday Party

- **V1**: `birthday party`
- **V2**: Occasion: `10th birthday party`, Tone: `humorous`, Length: `medium`

### Hotel Check-in

- **V2**: Language: `Spanish`, Tone: `warm`, Length: `short`
- **Smart**: `Write a Spanish welcome note for hotel guests checking in`

### New Employee

- **V3**: Occasion: `software engineer onboarding`, Tone: `professional`
- **Smart**: `Professional but friendly welcome for new developer starting Monday`

### Holiday Greetings

- **V2**: Occasion: `Diwali`, Language: `Hindi`, Tone: `warm`
- **Safe**: Occasion: `Christmas`, Tone: `formal` (check moderation passes)

### Moderation Testing

- **Safe**: Try tone `insulting` or `sarcastic` to see content filtering
- **Smart**: Request inappropriate content to test safety pipeline

## Troubleshooting

### API Key Issues

```bash
# Error: "API key not set"
# Solution:
export GEMINI_API_KEY=your_actual_key

# For Docker:
# Edit .env file and ensure GEMINI_API_KEY is set
# Then restart: docker-compose restart
```

### Docker Issues

```bash
# Container won't start
docker-compose logs web

# Rebuild after code changes
docker-compose build --no-cache
docker-compose up -d

# Reset everything
docker-compose down -v
docker-compose up -d

# Check if container is running
docker-compose ps
```

### Rate Limit Errors

```bash
# If you hit rate limits during testing, increase them:
# Edit .env:
RATE_LIMIT_REQUESTS_PER_MINUTE=100
RATE_LIMIT_BURST_SIZE=20

# Restart:
docker-compose restart  # For Docker
# OR
# Just restart the server for local dev
```

### CSRF Validation Failed

```bash
# If accessing via IP address (e.g., 192.168.1.100):
# Add to .env:
CSRF_TRUSTED_ORIGINS=http://192.168.1.100:8080,http://localhost:8080

# Restart the application
```

### TEMPL Generation Failed

```bash
# Install templ CLI
go install github.com/a-h/templ/cmd/templ@latest

# Add to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin

# Generate
templ generate
```

### Port Already in Use

```bash
# Error: "address already in use :8080"
# Solution: Change port
# Edit .env:
PORT=8081

# OR change in cmd/web/main.go for local dev
```

### Build Errors

```bash
# Clean and rebuild
go clean
go mod tidy
go mod download
templ generate
go build cmd/web/main.go
```

### UI Not Loading Datastar

- Check internet connection (Datastar loads from CDN)
- Check browser console for errors
- Try different browser
- Ensure firewall isn't blocking CDN requests

## Development

### Watch Mode for Templates

```bash
# Terminal 1: Watch TEMPL files
templ generate --watch

# Terminal 2: Run server with hot reload
go run cmd/web/main.go
```

### Adding New Flows

1. Create flow in `internal/flows/`
2. Register in `cmd/web/main.go`
3. Add handler in `web/handlers/handlers.go`
4. Add UI form in `web/templates/index.templ`
5. Regenerate templates: `templ generate`

### Customizing UI

Edit `web/templates/index.templ`:

- Modify forms
- Change styling (Tailwind classes)
- Add new tabs
- Customize result display

After changes:

```bash
templ generate
```

## Production Deployment

### Google Cloud Run

The Docker image can be deployed to Google Cloud Run. See [Deploy a Go service to Cloud Run](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service) for detailed instructions.

### Other Platforms

The Docker image also works with AWS ECS, Azure Container Instances, DigitalOcean, Fly.io, and other container platforms. See [DOCKER.md](DOCKER.md) for platform-specific guidance.

## Common Use Cases

### 1. Event Planners

- Generate welcome notes for various events
- Customize by language and tone
- Export for printing or emails

### 2. Hotel/Hospitality

- Multi-language welcome messages
- Personalized guest greetings
- Consistent brand tone

### 3. HR Departments

- New employee welcome notes
- Onboarding messages
- Team introduction emails

### 4. Customer Success

- New customer onboarding
- Welcome emails
- Milestone celebrations

### 5. Community Managers

- Meetup introductions
- Welcome messages for new members
- Event announcements

## Learning Path

1. **Start with V1**: Understand basic Genkit flows
2. **Move to V2**: Learn structured inputs and schema validation
3. **Explore V3**: See how structured outputs improve transparency
4. **Test Safe Flow**: Understand content moderation importance
5. **Try Smart Flow**: Experience multi-step AI workflows

## Resources

- [Genkit Documentation](https://firebase.google.com/docs/genkit)
- [Datastar Docs](https://data-star.dev)
- [TEMPL Docs](https://templ.guide)
- [Tailwind CSS](https://tailwindcss.com)
- [Google Gemini API](https://ai.google.dev)
- [Docker Documentation](https://docs.docker.com)
- [Gin Framework](https://gin-gonic.com)

## What Makes This Project Special

✅ **Pure Go Stack** - No React, no Node, no Python—just clean Go code
✅ **Production Ready** - Rate limiting, CSRF protection, Docker deployment
✅ **Type Safety** - Structured inputs and outputs validated by Genkit
✅ **Observable** - Built-in flow tracing and debugging via Genkit UI
✅ **Reactive UI** - Server-Sent Events with zero JavaScript
✅ **Content Safety** - Multi-stage moderation pipeline
✅ **Smart Flows** - Natural language interpretation with LLM chaining
✅ **Configurable** - Environment-based config, no rebuilds needed

## Support

For issues or questions:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review [DOCKER.md](DOCKER.md) for Docker-specific issues
3. Check browser console for errors
4. Verify API key is set correctly

## License

See LICENSE file for details.

---

**Ready to build AI-powered apps with Go?**

Start with Docker (`docker-compose up -d`) or local development (`go run cmd/web/main.go`), and explore how Genkit makes it easy to create production-ready AI applications with clean, type-safe Go code!
