# Epic Games Free Games API

A Go-based REST API that provides information about current and upcoming free games from the Epic Games Store.

## Features

-   🎮 Get current free games
-   📅 Get upcoming free games
-   🔄 Real-time data from Epic Games Store
-   📝 Detailed game information including:
    -   Title and description
    -   Original price and current status
    -   Promotional dates
    -   Game images
    -   Store links
    -   Seller information

## Tech Stack

-   Go 1.21+
-   Standard Go HTTP package for server
-   Clean Architecture pattern

## Installation

1. Clone the repository:

```bash
git clone https://github.com/RehanDias/free-games-api.git
cd free-games-api
```

2. Install dependencies:

```bash
go mod download
```

## Usage

### Running the Server

You can run the server using either of these commands:

```bash
# Using the API entry point
go run cmd/api/main.go

# Or using the server entry point
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### API Endpoints

#### GET /api/free-games

Returns both current and upcoming free games from Epic Games Store.

Example response:

```json
{
    "success": true,
    "timestamp": "2025-05-05T10:00:00Z",
    "data": {
        "current": [
            {
                "title": "Game Title",
                "description": "Game description",
                "status": "FREE NOW",
                "offerType": "BASE_GAME",
                "effectiveDate": "2025-05-05 10:00:00",
                "seller": "Game Studio",
                "price": {
                    "originalPrice": 19.99,
                    "formattedOriginalPrice": "$19.99",
                    "discount": "100%",
                    "current": "FREE"
                },
                "images": {
                    "wide": "https://...",
                    "thumbnail": "https://..."
                },
                "urls": {
                    "product": "https://store.epicgames.com/en-US/p/game-slug"
                },
                "availability": {
                    "endDate": "2025-05-12 10:00:00"
                }
            }
        ],
        "upcoming": [
            // Similar structure to current games
        ]
    }
}
```

## Project Structure

```
.
├── cmd/                    # Application entry points
│   ├── api/               # API server entry point
│   └── server/            # Alternative server entry point
├── internal/              # Private application code
│   ├── handlers/          # HTTP request handlers
│   ├── models/            # Data models
│   ├── server/            # Server configuration
│   ├── services/          # Business logic
│   └── utils/             # Utility functions
└── go.mod                 # Go module file
```

## Error Handling

The API returns standardized error responses:

```json
{
    "success": false,
    "timestamp": "2025-05-05T10:00:00Z",
    "error": {
        "message": "Error description",
        "code": 500
    }
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

-   Epic Games Store for providing the game data
-   The Go community for the amazing standard library
