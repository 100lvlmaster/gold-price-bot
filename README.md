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
go build -o gold-price-api ./cmd/api/main.go
```

## 3. How to Deploy using GitHub Actions

1. **Prerequisites on Target Server (DigitalOcean Droplet):**
   - Ensure target directory `/var/www/goldprice-bot` exists.
   - Copy the service file and enable it:
     ```bash
     sudo cp services/goldprice-bot.service /etc/systemd/system/
     sudo systemctl daemon-reload
     sudo systemctl enable goldprice-bot.service
     sudo systemctl start goldprice-bot.service
     ```

2. **Configure GitHub Repository Secrets:**
   Add the following secrets under **Settings > Secrets and variables > Actions**:
   - `DO_HOST`: Server IP address or domain.
   - `DO_USERNAME`: SSH login username.
   - `DO_SSH_KEY`: SSH private key.

3. **Deploy:**
   - Push to the `main` branch to trigger the workflow.
   - The workflow (`.github/workflows/deploy.yml`) builds the binary, copies it to the server, and restarts `goldprice-bot.service`.

## 4. Deploying with Docker Compose

1. **Configure Environment:**
   Ensure your `.env` file is created and set correctly (from `.env.example`).

2. **Start Services:**
   Run the following command to build the image and start the container in detached mode:
   ```bash
   docker compose up -d --build
   ```
   This will mount the SQLite database under `./data/database.db` to ensure persistence.

## 5. Running with Docker CLI

If you prefer using the Docker CLI directly instead of Docker Compose:

1. **Build the Image:**
   ```bash
   docker build -t gold-price-api .
   ```

2. **Run the Container:**
   ```bash
   docker run -d \
     --name goldprice-bot \
     -p 8080:8080 \
     --env-file .env \
     -e DB_PATH=/app/data/database.db \
     -v $(pwd)/data:/app/data \
     --restart unless-stopped \
     gold-price-api
   ```


