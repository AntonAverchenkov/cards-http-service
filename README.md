# cards-http-service

A simple demo of a stateful REST API service for a deck of cards.

## API

![api](/doc/api.png)

A few of the `POST` endpoints are duplicated as `GET` endpoints to simplify
testing in a browser. For example, it is possible to return a card by hitting
`GET /cards/return?card=qd` in a browser, rather than using `curl` (or similar).

## Session management

The service maintains a unique session for each browser client that connects to
it. The sessions are maintained by setting the `"session"` cookie. Each session
corresponds to a unique deck of cards view for the client. This is currently not
going to work if the cookies are not enabled.

### Session persistence

The sessions can be optionally persisted through server restarts:

```sh
./cards-http-service --sessions-persist-to path/to/file --sessions-restore-from path/to/file
```

The service will attempt to parse the `--sessions-restore-from` file on startup
and restore sessions from it. On shutdown, the service will write sessions to
the `--sessions-persist-to` file. This mechanism is not foolproof (it will not
work if the service hard-crashes). However, it should take care of most other
cases.

A valid sessions persistence file will look something like the one below
(`session-id serialized-deck-string`):

```
LnLgk_JPEZpRRtW9I5TUoM8M229EzcWTrmtz49YY4J4= thjhqhkhad2d3d
_yxvxLANbcXLPbPbKsPDZ2LLLS7gtzuozhQ0VYiLCZ8= 6c7c8c9ctcjcqckcah2h
```

## Install & run

```sh
go install github.com/AntonAverchenkov/cards-http-service
$GOPATH/bin/cards-http-service
```

Alternatively, the service can be run from docker:

```sh
docker build -t cards-http-service .
docker run -it -p 127.0.0.1:8080:8080/tcp cards-http-service
```

## Testing

### Unit tests

```sh
go test -v ./...
```

There are currently no environment-specific tests or tests that simulate client
requests (on my to-do list), however, it is possible to run the existing unit
tests from docker as well:

```sh
docker run -it -p 127.0.0.1:8080:8080/tcp cards-http-service go test -v ./...
```

### Testing in browser

Try the following endpoints in one or more browsers:

- http://localhost:8080/
- http://localhost:8080/cards
- http://localhost:8080/cards/shuffle
- http://localhost:8080/cards/deal
- http://localhost:8080/cards/return?card=ac

#### Short-form card encoding for /cards/return endpoint

The `/cards/return?card={card}` endpoing accepts both the long-form and the
short-form representations of a card. For example, to return `jack of spades`
to the deck either of the following will work:

- `/cards/return?card=jack%20of%20spades`
- `/cards/return?card=js`

Below is the full list of card suits and values, along with their short-form
representations:

| Suit     | Short |
|----------|-------|
| clubs    | `c`   |
| diamonds | `d`   |
| spades   | `s`   |
| hearts   | `h`   |


| Value | Short |
|-------|-------|
| ace   | `a`   |
| two   | `2`   |
| three | `3`   |
| four  | `4`   |
| five  | `5`   |
| six   | `6`   |
| seven | `7`   |
| eight | `8`   |
| nine  | `9`   |
| ten   | `t`   |
| jack  | `j`   |
| queen | `q`   |
| king  | `k`   |

