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

## 3. How to Deploy using GitHub Actions

1. **Prerequisites on Target Server (DigitalOcean Droplet):**
   - Ensure target directory `/var/www/goldprice-bot` exists:
     ```bash
     sudo mkdir -p /var/www/goldprice-bot
     ```
   - Place your `.env` file on the server (it is gitignored, so CI cannot ship it):
     ```bash
     sudo cp .env /var/www/goldprice-bot/.env
     ```
   - Copy the service file and enable it (one-time setup):
     ```bash
     sudo cp services/goldprice-bot.service /etc/systemd/system/
     sudo systemctl daemon-reload
     sudo systemctl enable --now goldprice-bot.service
     ```

2. **Configure GitHub Repository Secrets:**
   Add the following secrets under **Settings > Secrets and variables > Actions**:
   - `DO_HOST`: Server IP address or domain.
   - `DO_USERNAME`: SSH login username.
   - `DO_SSH_KEY`: SSH private key.

3. **Deploy:**
   - Push to the `main` branch to trigger the workflow.
   - The workflow (`.github/workflows/deploy.yml`) builds the binary, copies it to `/var/www/goldprice-bot`, copies the service file to `/etc/systemd/system`, and restarts `goldprice-bot.service`.

## 4. Migrating from Docker

If you previously ran the bot via Docker Compose, stop the old container to free port 8080:
```bash
cd /opt/goldprice-bot
docker compose down
```
