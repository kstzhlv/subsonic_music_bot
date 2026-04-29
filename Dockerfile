FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/subsonic-bot ./cmd

FROM gcr.io/distroless/static-debian12

COPY --from=build /out/subsonic-bot /subsonic-bot

ENTRYPOINT ["/subsonic-bot"]
