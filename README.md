# subsonic_music_bot

## Docker

1. Copy `.env.example` to `.env` and fill in your Subsonic credentials.
2. Build and run the bot:

```bash
docker compose up --build
```

The current entrypoint fetches the newest albums once and prints them to stdout.

### Private CA certificates

If your Subsonic server uses a certificate signed by a private or internal CA, mount that CA certificate into the container and set `SUBSONIC_CA_CERT_FILE` to its path inside the container.

Example:

```yaml
services:
  bot:
    build:
      context: .
    env_file:
      - .env
    volumes:
      - ./certs/subsonic-ca.pem:/certs/subsonic-ca.pem:ro
```
