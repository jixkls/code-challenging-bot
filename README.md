# 🤖 Code Challenger Bot

A Telegram bot that delivers Go programming challenges with a game-like progression system. Built in Go, deployed on [Fly.io](https://fly.io), with Supabase as the database backend.

**Bot:** [@CodeChallengerBot](https://t.me/CodeChallengerBot)
**Live URL:** [code-challenging-bot.fly.dev](https://code-challenging-bot.fly.dev/)

---

## About

Code Challenger Bot helps developers learn and practice Go through daily coding challenges. Users receive problems at different difficulty levels, track their progress with streaks, and level up as they solve more challenges.

### Features

- `/start` — Welcome message with personalized progress stats
- `/help` — List of available commands
- `/challenge` — Get a random challenge
- `/easy` — Get an easy Go challenge
- `/medium` — Get a medium Go challenge
- `/hard` — Get a hard Go challenge
- User progress tracking (level, challenges solved, streak)
- Webhook-based architecture for scalability

---

## Tech Stack

- **Language:** Go
- **Platform:** Telegram Bot API (webhooks)
- **Database:** Supabase (PostgreSQL) via GORM
- **Hosting:** Fly.io
- **Environment:** godotenv for local `.env` management

---

## Project Structure

```
telegram-bot/
├── main.go          # Webhook handler, command parsing, bot logic
├── go.mod
├── go.sum
├── .env             # Local environment variables (not committed)
├── Dockerfile       # Fly.io deployment config
└── fly.toml         # Fly.io app configuration
```

---

## Setup

### Prerequisites

- Go 1.21+
- A Telegram bot token from [BotFather](https://t.me/BotFather)
- Supabase project with database URL
- Fly.io CLI (`flyctl`) for deployment

## Roadmap

| Version | Description | Status |
|---------|-------------|--------|
| **v1** | Daily Go challenges with difficulty levels (easy/medium/hard) | ✅ Done |
| **v2** | Streak tracking, user progress, level progression (level N requires N challenges) | 🔧 In Progress |
| **v3** | External LLM API integration for dynamic challenge generation | 📋 Planned |
| **—** | Micro-SaaS launch with branding and monetization | 📋 Planned |

---

## Changelog

### [v0.4.0] - 2026-02-16
**Database & User Tracking**
- Integrated Supabase (PostgreSQL) as database backend via GORM
- Added `User` model with fields: telegram ID, username, level, challenges solved, current streak
- Implemented `GetOrCreateUser` for persistent user tracking
- `/start` command now shows personalized progress stats

### [v0.3.0] - 2026-02-15
**Challenge Commands**
- Added command parsing with `strings.HasPrefix` and switch statement
- Implemented `/challenge`, `/easy`, `/medium`, `/hard` commands
- Each difficulty delivers Go programming problems with code templates
- Telegram-formatted messages with emojis and code blocks
- Updated `/help` to list all available commands

### [v0.2.0] - 2026-02-15
**Production Deployment**
- Deployed bot to Fly.io with native Go buildpack
- Added health check endpoint (`/`) for Fly.io monitoring
- Dynamic port handling via `PORT` environment variable
- Migrated from ngrok local testing to production webhook
- Resolved TLS certificate verification issues in containerized environment
- Set Telegram webhook to production URL

### [v0.1.0] - 2026-02-11
**Webhook Foundation**
- Created Go structs for Telegram API: `User`, `Chat`, `Message`, `Update`
- Implemented `/webhook` HTTP handler to receive Telegram POST requests
- Built `sendReply` function to send messages back via Telegram API
- JSON marshaling/unmarshaling for Telegram message format
- Environment variable management with godotenv
- Local testing setup with ngrok
- Registered webhook with Telegram API
- Basic echo bot functionality (replies with "You said: ...")

### [v0.0.1] - 2026-02-07
**Project Kickoff**
- Defined project concept and bot roadmap (v1 → v2 → v3)
- Learned foundational Go concepts: structs, JSON tags, HTTP requests, `defer`, error handling
- Created `Challenge` struct with title, description, difficulty, example, and hint fields
- Created bot via BotFather

---

## Lessons Learned

Some key takeaways from building this project:

- **Error handling matters:** Always `return` after handling errors in Go to prevent nil pointer panics downstream.
- **TLS in containers:** Containerized environments may lack root CA certificates, causing HTTPS requests to fail. Install `ca-certificates` in your Dockerfile.
- **Webhooks > Polling:** Webhook architecture scales better than long polling for Telegram bots.
- **Exported fields:** Go struct fields must start with an uppercase letter to be visible to `encoding/json` and other packages. Use JSON tags to map to lowercase keys.

---

## License

This project is under development. License TBD.
