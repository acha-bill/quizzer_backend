# Quizzer API

Api for a quiz app.

## Docs

The doc definition is located in `docs/swagger.yaml` in the root of the repo.  
Go to https://editor.swagger.io/ and import the definition file.

## Running

1. install `go` and `mongodb`
2. In the root of the repo, create a `.env` file. See sample below
3. Run the app. `$ go run main.go`

## Docker
Coming soon

### Sample env
```env
MONGODB_URL=mongodb://localhost:27017
DATABASE_NAME=quizzer
JWT_SECRET=secret
DEBUGGING_ENABLED=false
```
