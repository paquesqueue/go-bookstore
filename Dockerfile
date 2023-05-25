FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test --tags=unit -v ./...

RUN go build -o ./out/go-bookstore .

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-bookstore /app/go-bookstore

CMD ["/app/go-bookstore"]