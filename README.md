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
      - COOKIE=
      - BARK_URL=
      # Use UTC Time
      - SCHEDULED_TIME=00:05
```

```bash
docker compose up -d
```