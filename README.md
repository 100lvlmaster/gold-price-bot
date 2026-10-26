# Gold Price Bot

A Go-based Telegram bot that fetches gold prices and sends updates/pings to registered users.

## 1. How to Install

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd goldprice-bot-go
   ```
2. **Install dependencies:**
   ```bash
   go mod download
   ```
3. **Configure environment:**
   Copy the example environment file and configure your variables:
   ```bash
   cp .env.example .env
   ```
   Open the `.env` file and set your `API_KEY` (Telegram bot token) and other configuration values.
4. **Create the SQLite database file:**
   Create an empty file matching your configured database path (defaults to `database.db`):
   ```bash
   touch database.db
   ```

## 2. How to Build

Build the production API binary:
```bash
CGO_ENABLED=1 go build -ldflags="-w -s" -o gold-price-api ./cmd/api/main.go
```

## 3. How to Deploy using GitHub Actions (self-hosted runner)

Deployment uses a **self-hosted runner** running directly on the VPS. No inbound SSH is required, so it works even when the provider blocks GitHub's datacenter IPs.

### 3.1 One-time server setup

1. **Install Go and build tools** (sqlite3 needs CGO/gcc):
   ```bash
   apt update
   apt install -y golang-go build-essential git
   ```

2. **Create the target directory and place `.env`** (gitignored, so CI cannot ship it):
   ```bash
   mkdir -p /var/www/goldprice-bot
   cp .env /var/www/goldprice-bot/.env
   ```

3. **Add a self-hosted runner:**
   - In the repo: **Settings → Actions → Runners → New self-hosted runner**.
   - Pick **Linux** and **x64**, then follow the **Download** and **Configure** commands on the server.
   - When prompted for a runner group, use the default. For labels, keep `self-hosted` (the workflow matches on it).

4. **Run the runner as root** so it can write to `/var/www` and `/etc/systemd/system`:
   ```bash
   sudo ./svc.sh install
   sudo ./svc.sh start
   ```

### 3.2 Deploy

Push to `main`. The workflow (`.github/workflows/deploy.yml`) runs on the server, builds the binary, installs it to `/var/www/goldprice-bot`, installs the systemd unit, and restarts `goldprice-bot.service`.

## 4. Migrating from Docker

If you previously ran the bot via Docker Compose, stop the old container to free port 8080:
```bash
cd /opt/goldprice-bot
docker compose down
```
