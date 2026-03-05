# 📡 PostPilot — LinkedIn Post Automation

> Automatically fetch trending **Tech News & AI** articles, generate LinkedIn posts with GPT-4, and publish them with one click from a beautiful approval dashboard.

![Go](https://img.shields.io/badge/Go-1.21-00acd7?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Status](https://img.shields.io/badge/Status-Active-brightgreen?style=flat-square)

---

## ✨ Features

- 🔍 **Auto-fetches** trending tech & AI news from 14 sources (Hacker News, TechCrunch, OpenAI Blog, Google AI, Wired, Ars Technica, dev.to, Reddit and more)
- 🤖 **GPT-4 writes** authentic LinkedIn posts in your voice — no corporate fluff
- 📋 **Approval dashboard** — review, edit, rewrite, approve or reject every post before it goes live
- ✏️ **Inline editing** — tweak posts directly in the browser
- ↺ **AI rewrite** — give an instruction like *"make it shorter"* or *"more technical"* and GPT rewrites it
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
│   │   └── generator.go              ← GPT-4 post generation & rewriting
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
- OpenAI API key → [platform.openai.com](https://platform.openai.com/api-keys)
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

OPENAI_API_KEY=sk-your-openai-key-here

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
| AI | OpenAI GPT-4 |
| Database | SQLite (GORM) |
| News | RSS feeds + NewsAPI |
| Publishing | LinkedIn UGC API |
| Frontend | Vanilla HTML/CSS/JS |
| Scheduler | robfig/cron |

---

## 🔒 Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `OPENAI_API_KEY` | ✅ Yes | GPT-4 post generation |
| `LINKEDIN_ACCESS_TOKEN` | ✅ Yes | Posting to LinkedIn |
| `LINKEDIN_PERSON_ID` | ✅ Yes | Your LinkedIn profile ID |
| `NEWS_API_KEY` | Optional | Extra news sources via NewsAPI |
| `AUTHOR_NAME` | Optional | Used in AI prompt (default: Dheeraj Reddy) |
| `PORT` | Optional | Server port (default: 8081) |

---

## 📌 Author

**Dheeraj Reddy** — Senior Software Engineer at Capital One
Go · Java · Python · AWS · Data Science

[![LinkedIn](https://img.shields.io/badge/LinkedIn-dheerajreddy-0077B5?style=flat-square&logo=linkedin)](https://www.linkedin.com/in/-dheerajreddy/)
[![GitHub](https://img.shields.io/badge/GitHub-dheerajreddy01-181717?style=flat-square&logo=github)](https://github.com/dheerajreddy01)

---

> ⚠️ **Note:** Use responsibly. Ensure posts comply with [LinkedIn's User Agreement](https://www.linkedin.com/legal/user-agreement). Review all AI-generated content before publishing.
