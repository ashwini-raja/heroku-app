# Heroku Go App

A simple Go web application ready for Heroku deployment.

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

- `main.go` - Main Go application
- `Procfile` - Tells Heroku how to run the app
- `go.mod` - Go module file
- `README.md` - This file

## Endpoints

- `/` - Main page
- `/health` - Health check endpoint 