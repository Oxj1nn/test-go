FROM golang:1.27.3-alpine3.24 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -o /app ./cmd/server

FROM scratch

COPY --from=build /app /app
USER 65532:65532

EXPOSE 8080
ENTRYPOINT ["/app"]
