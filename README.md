# Heroku Go App

A simple Go web application ready for Heroku deployment that makes authenticated requests to external services.

## Local Development

1. Make sure you have Go installed
2. Run the application locally:
   ```bash
   go run main.go
   ```
3. Visit http://localhost:8080

## Deploy to Heroku

1. Install Heroku CLI if you haven't already
2. Login to Heroku:
   ```bash
   heroku login
   ```

3. Create a new Heroku app:
   ```bash
   heroku create your-app-name
   ```

4. Deploy to Heroku:
   ```bash
   git add .
   git commit -m "Initial commit"
   git push heroku main
   ```

5. Open your app:
   ```bash
   heroku open
   ```

## Project Structure

- `main.go` - Main Go application with dynoid authentication
- `Procfile` - Tells Heroku how to run the app
- `go.mod` - Go module file with dependencies
- `README.md` - This file

## Endpoints

- `/` - Makes authenticated GET request to https://applink.staging.herokudev.com
- `/health` - Health check endpoint

## Authentication

The application uses the `github.com/heroku/x/dynoid` package to:
- Read the local dyno ID token
- Use it as a Bearer token in the Authorization header
- Authenticate requests to external services

## Dependencies

- `github.com/heroku/x/dynoid` - For reading dyno ID tokens 