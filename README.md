[![Go Report Card](https://goreportcard.com/badge/github.com/PonzaMatteo/router-open-data-hub)](https://goreportcard.com/report/github.com/PonzaMatteo/router-open-data-hub)

# Open Data Hub Challenge: Router to Different APIs

Repository for the Open Data Hub Router Challenge.

## Requirements

- Go: `1.21.4`
    - https://go.dev

## Commands

### Start Application

The following command will start a server listening to port `8080`:

```bash
go run ./cmd/main/main.go
```

To try it out in the browser head over to the following link:

- [Mobility Event Sample](http://localhost:8080/v2/flat,event/*/latest?limit=200&offset=0&where=evuuid.eq.53a6343f-e524-51ea-a280-4cc4c1bc7ff3&shownull=false&distinct=true)

### Run Test

```bash
go test ./...
```

## Useful Links

- [Tourism API Docs](https://tourism.opendatahub.com/swagger/index.html)
- [Mobility API DOCS](https://swagger.opendatahub.com/?url=https://mobility.api.opendatahub.com/v2/apispec#/Mobility%20V2)
