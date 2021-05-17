# cards-http-service

A simple demo of a stateful REST API service for a deck of cards.

## API

![api](/doc/api.png)

A few of the `POST` endpoints are duplicated as `GET` endpoints to simplify
testing in a browser. For example, it is possible to return a card by hitting
`GET /cards/return?card=3d` in a browser, rather than using `curl` (or similar)
with `json` payload.

## Session management

The service maintains a unique session for each browser that connects to it.
The sessions are maintained by setting the session cookie. Each session
corresponds to a single deck view. This is currently not going to work if
the cookies are not accepted by the browser.
