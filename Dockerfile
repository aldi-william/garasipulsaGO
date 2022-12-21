FROM golang:latest as build

RUN useradd -u 1001 go

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./main

FROM scratch

ENV GIN_MODE=release

WORKDIR /app

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/main /app/main
COPY --from=build /app/.env /app/.env

USER go

EXPOSE 8083

ENTRYPOINT ["/app/main"]
