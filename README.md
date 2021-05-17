# cards-http-service

A simple demo of a stateful REST API service for a deck of cards.

## API

![api](/doc/api.png)

A few of the `POST` endpoints are duplicated as `GET` endpoints to simplify
testing in a browser. For example, it is possible to return a card by hitting
`GET /cards/return?card=qd` in a browser, rather than using `curl` (or similar)
with `JSON` payload.

## Session management

The service maintains a unique session for each browser that connects to it.
The sessions are maintained by setting the session cookie. Each session
corresponds to a single deck view. This is currently not going to work if
the cookies are not enabled.

### Session persistence

The sessions can be optionally persisted through server restarts:

```sh
$ ./cards-http-service --sessions-persist-to path/to/file --sessions-restore-from path/to/file
```

If specified, the service will attempt to parse the given file on startup for
any persisted sessions. On shutdown, the server will write sessions to the
given file. This mechanism is not foolproof (it will not work if the service
hard-crashes). However, it should take care of most other cases.

## Installation

```sh
go install github.com/AntonAverchenkov/cards-http-service
```

## Testing

### Unit tests

```sh
go test -v --short ./...
```

### Testing in browser

Start up the service as follows:

```sh
$ ./cards-http-service --address :8080
```

Try the following in multiple browsers:

- http://localhost:8080/
- http://localhost:8080/cards
- http://localhost:8080/cards/shuffle
- http://localhost:8080/cards/deal
- http://localhost:8080/cards/return?card=ac

#### Short-form card encoding for /cards/return

`/cards/return?card={card}` accepts both the long-form and the short-form
representations of a card. For example, to return `jack of spades`
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

