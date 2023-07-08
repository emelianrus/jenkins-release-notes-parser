# Relase notes parser


## Run locally


```
docker-compose up -d - create DB
cd frontend && npm install && npm start - run WEB
go run . - run main app
```

later would be moved to compose


## Github personal token

by default github has limit to 60 request per hour

you can increase this limit with personal token

1. go to page https://github.com/settings/tokens (you should be autorized)

2. generate personal classic token with permission:
```
public_repo
```
3. set env var GITHUB_PERSONAL_TOKEN=<YOUR_TOKEN> or .env file

will take env var and if not set .env file else will use public api without token (60 requests per hour)