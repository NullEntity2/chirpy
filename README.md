# Chirpy
Chirpy is a simple Twitter-like clone. It has user authentication, authorization, creation of chirps and listing chirps. This was created as a practice program to learn about HTTP servers in Go from Boot.Dev.

## Installation
1. Install Go
2. Clone the repo locally
3. `go build`
4. `./chirpy`

## Configuration
Create a file named .env with the following contents:
```
DB_URL="postgres://bryce:@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="<jwt secret>"
POLKA_KEY="<webhook api key>"
```
Update the DB_URL to point to your database. ?sslmode=disable may be required.

## Endpoints
/app/ - Serves static content in /html
GET /api/healthz - Shallow health check
GET /admin/metrics - Hit count of /app/ static content pages
POST /admin/reset - Clears all users from the database for testing purposes
POST /api/users - Create a user
POST /api/chirps - Create a chirp
GET /api/chirps - List all chirps
GET /api/chirps/{id} - Get chirp by id
POST /api/login - Log in as a user, returns a short-llved access token and a long term refresh token
POST /api/refresh - Given a long term refresh token, returns a short term access token
POST /api/revoke - Revoke a given refresh token
PUT /api/users - Update the current user's username or password
DELETE /api/chirps/{chirpID} - Delete a chirp
POST /api/polka/webhooks - Webhook function to mark a user as Chirpy Red. Requires an API key.