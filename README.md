# Daebak Korean Backend

## Development
- docker-compose up
- docker exec -it daebak-web_db_1 psql -U onehappyfellow -d daebak
- create table if not exists
- go run main.go

## Testing
At some point, I will add unit and/or integration tests, but for now this checklist allows some level of confidence.
- create user
- profile page
- logout > profile should redirect
- login
- logout
- fail login
- reset password
- profile
- logout
- fail create (same email)

## TODO
- error toasts
- csrf protection
- password reset email
- account creation confirmation email
- migrations
- content page (article)
- user history
- delete all user history