FROM golang:1.16

ENV GOPROXY=https://proxy.golang.org

WORKDIR /go/src/cards-http-service
COPY . .

ENV ADDRESS=":8080"
ENV SESSIONS_PERSIST_TO="/tmp/cards-http-service.sessions"
ENV SESSIONS_RESTORE_FROM="/tmp/cards-http-service.sessions"

ENV GOPROXY=https://proxy.golang.org,direct

RUN go get -d -v ./...     && \
    go test -v ./...       && \
    go install -v ./...    && \
    go clean -i -modcache  && \
    go clean -i -cache     && \
    go clean -i -testcache

CMD ["cards-http-service"]
