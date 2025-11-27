# Docker Setup Guide

This guide explains how to build and run the Welcome Note Generator using Docker.

## Prerequisites

- Docker installed on your system
- Docker Compose (optional, but recommended)

## Quick Start

### 1. Obtain Google API Key

Get your Google AI API key from [Google AI Studio](https://aistudio.google.com/app/api-keys)

### 2. Generate CSRF Key

Generate a secure CSRF key:

```bash
openssl rand -hex 32
```

### 3. Create Environment File

Copy the example environment file and add your keys:

```bash
cp .env.example .env
# Edit .env and add:
# - Your GEMINI_API_KEY
# - Your generated CSRF_KEY
```

### 4. Build and Run with Docker Compose

```bash
# Build and start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

### 5. Access the Application

Open your browser and navigate to:

```
http://localhost:8080
```

## Manual Docker Commands

If you prefer to use Docker directly without Docker Compose:

### Build the Image

```bash
docker build -t welcome-note-generator:latest .
```

### Run the Container

```bash
docker run -d \
  --name welcome-note-generator \
  -p 8080:8080 \
  -e GEMINI_API_KEY=your_gemini_api_key_here \
  -e CSRF_KEY=your_64_character_hex_string_here \
  -e RATE_LIMIT_REQUESTS_PER_MINUTE=30 \
  -e RATE_LIMIT_BURST_SIZE=5 \
  welcome-note-generator:latest
```

### View Logs

```bash
docker logs -f welcome-note-generator
```

### Stop and Remove Container

```bash
docker stop welcome-note-generator
docker rm welcome-note-generator
```

## Environment Variables

### Required Variables

| Variable         | Description                            | Example                                                               |
| ---------------- | -------------------------------------- | --------------------------------------------------------------------- |
| `GEMINI_API_KEY` | Gemini AI API key for Genkit           | Get from [Google AI Studio](https://makersuite.google.com/app/apikey) |
| `CSRF_KEY`       | 32-byte hex string for CSRF protection | Generate with `openssl rand -hex 32`                                  |

### Optional Variables

| Variable                         | Description                                   | Default |
| -------------------------------- | --------------------------------------------- | ------- |
| `PORT`                           | Server port                                   | `8080`  |
| `CSRF_TRUSTED_ORIGINS`           | Comma-separated trusted origins (IPs/domains) | ` `     |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Max requests per minute per IP                | `30`    |
| `RATE_LIMIT_BURST_SIZE`          | Burst size for rate limiter                   | `5`     |
| `RATE_LIMIT_CLEANUP_INTERVAL`    | Cleanup interval (Go duration)                | `5m`    |
| `RATE_LIMIT_LIMITER_TTL`         | Limiter TTL (Go duration)                     | `15m`   |

## Configuration Changes

To adjust rate limiting without rebuilding:

1. Stop the container
2. Update environment variables in `.env` or `docker-compose.yml`
3. Restart the container

```bash
docker-compose down
# Edit .env or docker-compose.yml
docker-compose up -d
```

## Health Check

The container includes a health check that verifies the application is running:

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' welcome-note-generator

# Or with docker-compose
docker-compose ps
```

## Image Details

- **Base Image**: Alpine Linux (latest)
- **Go Version**: 1.25
- **Binary Size**: Optimized with stripped symbols
- **Security**: Runs as non-root user (uid 1000)
- **Runtime Dependencies**: ca-certificates, tzdata

## Troubleshooting

### Container fails to start

Check the logs:

```bash
docker-compose logs
```

### Missing API key errors

Ensure both `GEMINI_API_KEY` and `CSRF_KEY` are set:

```bash
# Check if environment variables are set
docker exec welcome-note-generator env | grep -E "GEMINI_API_KEY|CSRF_KEY"
```

### CSRF errors

Ensure `CSRF_KEY` is set and is a valid 64-character hex string:

```bash
# Verify your CSRF_KEY
echo $CSRF_KEY | wc -c  # Should output 65 (64 chars + newline)
```

### Port already in use

Change the port mapping in `docker-compose.yml`:

```yaml
ports:
  - "3000:8080" # Map host port 3000 to container port 8080
```

### Rate limiting too strict/lenient

Adjust the rate limit settings in `.env`:

```bash
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_BURST_SIZE=10
```

Then restart the container:

```bash
docker-compose restart
```
