# Epic Games Free Games API

This is a Go API that fetches current and upcoming free games from the Epic Games Store.

## Local Development

1. Make sure you have Go installed (1.16 or later)
2. Clone the repository
3. Install dependencies:

```bash
go mod download
```

4. Run the server:

```bash
go run cmd/api/main.go
```

The server will start at `http://localhost:3000`

## API Endpoints

-   `GET /free-games` - Returns current and upcoming free games

## Deployment

This project can be deployed to Vercel:

1. Install Vercel CLI:

```bash
npm i -g vercel
```

2. Login to Vercel:

```bash
vercel login
```

3. Deploy:

```bash
vercel
```

## Structure

```
├── cmd/
│   └── api/
│       └── main.go       # Application entry point
├── internal/
│   ├── handlers/         # HTTP request handlers
│   ├── models/           # Data models
│   └── services/         # Business logic
├── vercel.json          # Vercel deployment config
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```
