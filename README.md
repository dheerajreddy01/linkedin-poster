# 📡 PostPilot — LinkedIn Post Automation

> Automatically fetch trending **Tech News & AI** articles, generate LinkedIn posts with Groq (Llama 3.3), and publish them with one click from a beautiful approval dashboard.

![Go](https://img.shields.io/badge/Go-1.21-00acd7?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Status](https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square)

---

## ✨ Features

- 🔍 **Auto-fetches** trending tech & AI news from 14 sources (Hacker News, TechCrunch, OpenAI Blog, Google AI, Wired, Ars Technica, dev.to, Reddit and more)
- 🤖 **Groq (Llama 3.3) writes** authentic LinkedIn posts in your voice — no corporate fluff, fast + cheap generation
- 📋 **Approval dashboard** — review, edit, rewrite, approve or reject every post before it goes live
- ✏️ **Inline editing** — tweak posts directly in the browser
- ↺ **AI rewrite** — give an instruction like *"make it shorter"* or *"more technical"* and the model rewrites it
- 🚀 **One-click publish** to LinkedIn via official API
- ⏰ **Auto-scheduler** — fetches new articles and generates drafts every 6 hours
- 🗂️ **Topic filters** — filter by AI & LLMs or Tech News
- 💾 **SQLite** — zero setup database, everything stored locally

---

## 🗂 Project Structure

```
linkedin-poster/
├── .env                              ← Your API keys (never commit this)
├── go.mod                            ← Go module dependencies
├── Dockerfile                        ← Container build for deployment (Fly.io etc.)
├── fly.toml                          ← Fly.io app config (persistent volume + service)
├── start.sh                          ← Mac/Linux startup script
├── start.ps1                         ← Windows startup script
│
├── cmd/
│   └── server/
│       └── main.go                   ← Entry point, routes, server
│
├── internal/
│   ├── api/
│   │   └── handlers/
│   │       └── handlers.go           ← REST API endpoints
│   ├── ai/
│   │   └── generator.go              ← Groq (Llama 3.3) post generation & rewriting
│   ├── news/
│   │   └── fetcher.go                ← RSS + NewsAPI fetcher (Tech & AI)
│   ├── linkedin/
│   │   └── client.go                 ← LinkedIn UGC API poster
│   ├── scheduler/
│   │   └── scheduler.go              ← Cron: fetch every 6hrs, auto-generate
│   ├── db/
│   │   └── db.go                     ← SQLite init & helpers
│   └── models/
│       └── models.go                 ← Post, NewsItem, Settings structs
│
├── frontend/
│   └── index.html                    ← Full approval dashboard UI
│
└── data/
    └── poster.db                     ← SQLite database (auto-created)
```

---

## ⚙️ Setup

### Prerequisites
- [Go 1.21+](https://go.dev/dl/)
- Groq API key *(free)* → [console.groq.com/keys](https://console.groq.com/keys)
- LinkedIn Developer App → [linkedin.com/developers/apps](https://www.linkedin.com/developers/apps)
- NewsAPI key *(optional, free)* → [newsapi.org](https://newsapi.org/register)

### 1 — Clone the repo

```bash
git clone https://github.com/dheerajreddy01/linkedin-poster.git
cd linkedin-poster
```

### 2 — Configure `.env`

```env
PORT=8081
DB_PATH=./data/poster.db

GROQ_API_KEY=gsk-your-groq-key-here
GROQ_MODEL=llama-3.3-70b-versatile        # optional, this is the default

NEWS_API_KEY=your-newsapi-key-here        # optional but recommended

LINKEDIN_ACCESS_TOKEN=your-access-token
LINKEDIN_PERSON_ID=your-person-id

AUTHOR_NAME=Dheeraj Reddy
```

#### Getting LinkedIn credentials
1. Go to [linkedin.com/developers/apps](https://www.linkedin.com/developers/apps) → Create App
2. Under **Products**, request access to **Share on LinkedIn**
3. Under **Auth**, generate an **Access Token** with scope `w_member_social`
4. Your **Person ID** is the identifier in your LinkedIn profile URL

### 3 — Start

**Mac / Linux:**
```bash
chmod +x start.sh
./start.sh
```

**Windows:**
```powershell
.\start.ps1
```

### 4 — Open the dashboard

```
http://localhost:8081
```

---

## ☁️ Deploy (Fly.io)

This app is a stateful Go server (SQLite + a 6-hourly cron job), so it needs an always-on host with
persistent disk — not a static/serverless host. [Fly.io](https://fly.io) fits well: it runs the Go
binary directly and gives you a small persistent volume for `data/poster.db`.

### 1 — Install flyctl and log in

```bash
curl -L https://fly.io/install.sh | sh
fly auth login
```

### 2 — Launch the app

From the repo root (`Dockerfile` and `fly.toml` are already included):

```bash
fly launch --no-deploy   # pick a unique app name; reuses the existing fly.toml
fly volumes create poster_data --size 1 --region <same region as fly.toml>
```

### 3 — Set secrets (never put these in fly.toml)

```bash
fly secrets set GROQ_API_KEY=gsk-your-groq-key-here
fly secrets set LINKEDIN_ACCESS_TOKEN=your-access-token
fly secrets set LINKEDIN_PERSON_ID=your-person-id
fly secrets set NEWS_API_KEY=your-newsapi-key-here   # optional
```

### 4 — Deploy

```bash
fly deploy
fly open
```

`fly.toml` mounts the volume at `/data` and sets `DB_PATH=/data/poster.db`, so drafts and settings
survive redeploys and restarts.

---

## 🖥️ Dashboard

| Action | Description |
|--------|-------------|
| **⚡ Fetch New Posts** | Manually trigger news fetch + AI generation |
| **✓ Approve** | Mark post as approved and ready to publish |
| **✕ Reject** | Skip this post |
| **✎ Edit** | Edit post content inline in browser |
| **↺ Rewrite** | Give AI an instruction to rewrite the post |
| **🚀 Post Now** | Publish approved post directly to LinkedIn |

---

## 📰 News Sources

### AI & LLMs
- OpenAI Blog
- Google AI Blog
- Hugging Face Blog
- MIT AI News
- r/MachineLearning
- r/artificial

### Tech News
- Hacker News
- TechCrunch
- Ars Technica
- The Verge
- Wired
- dev.to
- InfoQ
- r/technology

---

## 🔌 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/posts` | Get all posts (filter by `?status=draft`) |
| `GET` | `/api/posts/stats` | Get post counts by status |
| `PUT` | `/api/posts/:id/approve` | Approve a post |
| `PUT` | `/api/posts/:id/reject` | Reject a post |
| `PUT` | `/api/posts/:id/edit` | Update post content |
| `POST` | `/api/posts/:id/regenerate` | Rewrite with AI instruction |
| `POST` | `/api/posts/:id/post` | Publish to LinkedIn |
| `GET` | `/api/settings` | Get settings |
| `PUT` | `/api/settings` | Update settings |

---

## ⏰ Scheduler

| Job | Schedule | Action |
|-----|----------|--------|
| Fetch & Generate | Every 6 hours | Fetches news → generates post drafts |
| Startup fetch | On launch (after 3s) | Immediate first fetch |

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go + Gin |
| AI | Groq (Llama 3.3, OpenAI-compatible API) |
| Database | SQLite (GORM) |
| News | RSS feeds + NewsAPI |
| Publishing | LinkedIn UGC API |
| Frontend | Vanilla HTML/CSS/JS |
| Scheduler | robfig/cron |
| Hosting | Fly.io (Docker + persistent volume) |

---

## 🔒 Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GROQ_API_KEY` | ✅ Yes | Groq post generation |
| `LINKEDIN_ACCESS_TOKEN` | ✅ Yes | Posting to LinkedIn |
| `LINKEDIN_PERSON_ID` | ✅ Yes | Your LinkedIn profile ID |
| `NEWS_API_KEY` | Optional | Extra news sources via NewsAPI |
| `GROQ_MODEL` | Optional | Groq model id (default: `llama-3.3-70b-versatile`) |
| `AUTHOR_NAME` | Optional | Used in AI prompt (default: Dheeraj Reddy) |
| `PORT` | Optional | Server port (default: 8081) |
| `DB_PATH` | Optional | SQLite file path (default: `./data/poster.db`) |

---

## 📌 Author

**Dheeraj Reddy** — Senior Software Engineer at Capital One
Go · Java · Python · AWS · Data Science

[![LinkedIn](https://img.shields.io/badge/LinkedIn-dheerajreddy-0077B5?style=flat-square&logo=linkedin)](https://www.linkedin.com/in/-dheerajreddy/)
[![GitHub](https://img.shields.io/badge/GitHub-dheerajreddy01-181717?style=flat-square&logo=github)](https://github.com/dheerajreddy01)

---

> ⚠️ **Note:** Use responsibly. Ensure posts comply with [LinkedIn's User Agreement](https://www.linkedin.com/legal/user-agreement). Review all AI-generated content before publishing.
