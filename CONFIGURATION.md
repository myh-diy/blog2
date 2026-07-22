# Configuration

The application reads `config.yml` with Viper. Configuration precedence is:

1. Environment variables
2. `config.yml`
3. Built-in defaults

Set `CONFIG_FILE` to load a YAML file from another path.

## Production

Keep non-sensitive settings in `config.yml` and inject secrets with Compose:

```yaml
services:
  blog:
    image: ghcr.io/myh-diy/blog:latest
    expose:
      - "8080"
    volumes:
      - blog-uploads:/app/uploads
      - blog-data:/app/data
    environment:
      JWT_SECRET: replace-with-a-long-random-secret
      ADMIN_USERNAME: admin
      ADMIN_PASSWORD: replace-with-a-strong-password
      PUBLIC_URL: https://blog.example.com
      GIN_MODE: release
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://127.0.0.1:8080/api/health"]
      interval: 10s
      timeout: 3s
      retries: 5
      start_period: 10s
```

`ADMIN_USERNAME` and `ADMIN_PASSWORD` are used only when the first user is
created. Changing them later does not reset an existing account.

Supported legacy environment variables such as `PORT`, `DB_PATH`,
`JWT_SECRET`, `GIN_MODE`, and `EXPORTER_METRICS_URL` remain valid.
