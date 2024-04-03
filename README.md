# Decks API

An API to create and shuffle decks, open them, and draw cards.

## Getting started

You will need Go version 1.22.1 or newer installed. Check https://go.dev/doc/install for details.

Then simply extract the code from the source provided, or clone from https://github.com/rnkjnk/decks-api.

Once you have the source, make sure that required modules are installed:
```
go mod download
```

This will make sure the following dependencies are installed:
```
github.com/gin-gonic/gin v1.9.1
github.com/google/uuid v1.6.0
gopkg.in/yaml.v3 v3.0.1
```
The port at which the API will run is configurable in `config.yaml`, the default is 8080. Make sure no firewall is blocking you. SSL is not supported.

Run the API by running the only entry point:
```
go run cmd/main.go
```

## Usage

The following routes will be available:

### Create deck
```
POST   /deck
```
Creates and returns a new deck.

Query parameters: 
`shuffle` Boolean. if set to true, the deck will be shuffled. Default is true. Example: `POST /deck
`cards` Comma separated list of card codes which can be used to choose from which cards to create a deck. If ommited, all cards are used. Order of the supplied list is irrelevant, if not shuffled, the order provided in the config file will be used. Cards that don't exist in the config (such as 0C) will be ignored. Badly formatted card codes (such as 10C) will cause an error.

Example: `POST /deck?shuffle=false&cards=AC,AH,AD,AS`

Return value:
```
{
    "deck_id": "02b1ea53-4785-4f74-b0fc-90c4b945de12",
    "shuffled": false,
    "remaining": 4
}
```

### Open deck
```
GET    /deck/:id/open
```
Opens a deck.

URL parameters: 
`:id` The ID (uuid) of the deck requested. This parameter is mandatory.

Example: `deck/133316bd-1cb4-4b57-af75-43bd54fe60cd/open`

Return value:
```
{
    "deck_id": "02b1ea53-4785-4f74-b0fc-90c4b945de12",
    "shuffled": true,
    "remaining": 52,
    "cards": [
        {
            "suit": "DIAMONDS",
            "value": "2",
            "code": "2D"
        },
        ...
    ]
}
```

### Draw cards
```
POST   /deck/:id/draw-cards
```
Draws the given number of cards from the deck. 

URL parameters: 
`:id` The ID (uuid) of the deck requested. This parameter is mandatory.

Query parameters: 
`draw` uint8. The number of cards to draw.

Example: `deck/133316bd-1cb4-4b57-af75-43bd54fe60cd/draw-cards?draw=2`

Return value:
```
{
    "cards": [
        {
            "suit": "DIAMONDS",
            "value": "2",
            "code": "2D"
        },
        {
            "suit": "HEARTS",
            "value": "QUEEN",
            "code": "QH"
        }
    ]
}
```

## Configuration

The `config.yaml` file offers configurations for the API (port number), and for the decks service (names of card values and suits).

The card values and suits by default will generate a standard 52-card deck without jokers, with aces being first (lowest). The configuration can be modified (including the order of cards) with the following limitations:
- All values and suits respoectively must start with unique letters, since those are used to generate two-character card codes.
- The codes constructed from the first letters of these values are case sensitive.
- The total number of possible cards (so, number of suits times number of values) cannot exceed 255

## Tests
All unit tests are in the `tests` subdirectory. To run tests, simply run:
```
go test ./tests/services
```

Unit tests cover only the core service directly, and the in-memory store indirectly, due to time limitation. Methods related to loading configuration from yaml, and cofiguring the API routes, are not tested.

For any additional questions, please feel free to contact me at rnkjnk@gmail.com