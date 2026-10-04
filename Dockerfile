FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/potager ./cmd/potager

FROM golang:1.25-alpine AS goose
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

# Applies the migrations; run by the `migrate` service in compose.yaml.
FROM alpine:3.22 AS migrate
COPY --from=goose /go/bin/goose /usr/local/bin/goose
COPY migrations /migrations
ENTRYPOINT ["goose"]

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/potager /potager
COPY migrations /migrations
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/potager"]
