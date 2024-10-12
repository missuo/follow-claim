[![GitHub Workflow][1]](https://github.com/missuo/follow-claim/actions)
[![Go Version][2]](https://github.com/missuo/follow-claim/blob/main/go.mod)
[![Docker Pulls][3]](https://hub.docker.com/r/missuo/follow-claim)

[1]: https://img.shields.io/github/actions/workflow/status/missuo/follow-claim/docker.yaml?logo=github
[2]: https://img.shields.io/github/go-mod/go-version/missuo/follow-claim?logo=go
[3]: https://img.shields.io/docker/pulls/missuo/follow-claim?logo=docker

# Follow Claim

Follow Claim is a simple tool that uses a cron job to claim daily rewards from the Follow app.

## Usage

### Docker
```bash
docker run -d --name follow-claim -e COOKIE="your cookie" -e BARK_URL="your bark url" -e SCHEDULED_TIME="00:05" missuo/follow-claim
```

### Docker Compose

```
mkdir follow-claim && cd follow-claim
nano compose.yaml
```

```yaml
services:
  follow-claim:
    container_name: follow-claim
    image: missuo/follow-claim:latest
    restart: unless-stopped
    environment:
      # Cookie (Support multiple cookies, separated by commas) (Required)
      - COOKIE=
      # Bark URL (Optional)
      - BARK_URL=
      # Use UTC Time (UTC 00:05 is 08:05 in China) (Optional, Default: 00:05)
      - SCHEDULED_TIME=00:05
      # Telegram Bot Token (Optional)
      - TELEGRAM_BOT_TOKEN=
      # Telegram Chat ID (Optional)
      - TELEGRAM_CHAT_ID=
```

```bash
docker compose up -d
```
