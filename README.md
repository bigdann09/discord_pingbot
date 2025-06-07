## Network Ping Bot

### Purpose
The Network Ping Bot is a Go application that monitors server availability by pinging a list of hostnames stored in a PostgreSQL database using GORM. Users can add hostnames with `!add <hostname>`, and check their status with `!status <hostname>` via Discord commands. The bot continuously monitors all persisted hostnames every 4 seconds and sends alerts to a dedicated Discord alert channel if any are down. The bot is Dockerized and uses Docker Compose to connect with PostgreSQL, and it can be deployed on Render for free hosting. A task list for this and related projects is available for download below.

## Setup
### Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- PostgreSQL (local or hosted)
- Discord account and bot token
- Render account for deployment

### Local Setup

1. ### Clone the repository:
```bash
git clone https://github.com/bigdann09/discord_pingbot.git
cd network-ping-bot
```

2. ### Install dependencies:
```bash
go mod init pingbot
go get github.com/bwmarrin/discordgo
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

3. ### Set up Discord bot:

Create a Discord bot at Discord Developer Portal.
Copy the bot token from the "Bot" tab (not Client Secret).
In the "OAuth2 > URL Generator", select the bot scope and grant permissions: Send Messages, Read Messages/View Channels, and Read Message History.
Use the generated URL to invite the bot to your server.
Enable Developer Mode in Discord (User Settings > Appearance > Developer Mode), right-click the command channel and alert channel, and copy their IDs.


4. ### Configure environment variables:Create a .env file in the project root:
```bash
DISCORD_TOKEN=your-discord-bot-token
DISCORD_CHANNEL_ID=your-command-channel-id
DISCORD_ALERT_CHANNEL_ID=your-alert-channel-id
POSTGRES_DSN=host=postgres port=5432 user=postgres password=your-password dbname=pingbot sslmode=disable
```

5. ### Set up Docker:

Ensure Docker and Docker Compose are installed.
Create a docker-compose.yml file (see below) to run the bot and PostgreSQL.
Start the services:docker-compose up -d




6. ### Using the bot:

- In the Discord command channel:
    - Add a hostname: `!add example.com`
    - Remove a hostname: `!removeserver example.com`
    - Check status: `!status example.com`


The bot monitors all hostnames every 4 seconds and sends alerts to the alert channel if any are down.



Docker Setup
Use the following Dockerfile and docker-compose.yml to run the bot with PostgreSQL.
Dockerfile
```dockerfile
FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o pingbot
CMD ["./pingbot"]
```

docker-compose.yml
```yml
services:
    bot:
        container_name: pingbot
        build:
        context: .
        dockerfile: Dockerfile
        environment:
        - DISCORD_TOKEN=${DISCORD_TOKEN}
        - DISCORD_CHANNEL_ID=${DISCORD_CHANNEL_ID}
        - DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
        networks:
        - bot_network
        depends_on:
        - postgresdb

    postgresdb:
        container_name: postgresdb
        image: postgres:latest
        environment:
        POSTGRES_USER: ${POSTGRES_USER}
        POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
        POSTGRES_DB: ${POSTGRES_DB}
        ports:
        - ":5432"
        volumes:
        - postgres_data:/var/lib/postgresql/data
        networks:
        - bot_network
        healthcheck:
        test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
        interval: 5s
        timeout: 5s
        retries: 5

networks:
    bot_network:
        driver: bridge

volumes:
    postgres_data:
```


Create a `.env` file with POSTGRES_PASSWORD (e.g., POSTGRES_PASSWORD=your-password).
Run `docker-compose up -d` to start the bot and PostgreSQL.

Deployment on Render

1. ### Create a Render account at render.com.
2. ### Set up PostgreSQL on Render:
    - Create a new PostgreSQL database in Render.
    - Copy the internal connection string (e.g., `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable`).


3. ### Create a new Web Service:
Select "Docker" as the runtime.
Use the same GitHub repository containing Dockerfile and docker-compose.yml.


4. ### Configure environment variables in Render's dashboard:
Add DISCORD_TOKEN, DISCORD_CHANNEL_ID, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB (use the internal connection string from Render's PostgreSQL).


5. ### Deploy:
Push your code to a GitHub repository.
Connect the repository to Render and deploy.


6. ### Interact with the bot:
Use Discord commands (`!add`, `!status`) in the command channel.
Monitor alerts in the alert channel.


### Demo
A live demo is not hosted due to the need for Discord bot token and channel configuration. Follow the setup instructions to run locally with Docker or deploy on Render to interact with the bot.
Notes

The bot pings hostnames every 4 seconds and sends alerts to the alert channel.
Ensure the Discord bot has permissions to read and send messages in both the command and alert channels.
Hostnames must be unique in the database.
