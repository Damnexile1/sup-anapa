# Agent Commands

## Build and Deployment

```bash
# Full rebuild and restart
docker compose down && docker compose build && docker compose up -d

# Check logs
docker compose logs app --tail=20

# Check status
docker compose ps
```

## Database Management

```bash
# Connect to database
docker compose exec postgres psql -U postgres -d sup_anapa

# Check tables
docker compose exec postgres psql -U postgres -d sup_anapa -c "\dt"

# Check data in a table
docker compose exec postgres psql -U postgres -d sup_anapa -c "SELECT * FROM admins;"
```

## Application Testing

```bash
# Test API endpoints
curl -s http://localhost:8080/
curl -s -c cookies.txt -X POST http://localhost:8080/admin/login -d "username=admin&password=admin123"
curl -s -b cookies.txt http://localhost:8080/admin
```

## Troubleshooting

```bash
# Check container logs for errors
docker compose logs app | grep -i error

# Check if port is accessible
curl -I http://localhost:8080

# Restart specific service
docker compose restart app
```