FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate

RUN go build -o /app/rtc -tags netgo ./cmd

RUN chown 1000:1000 /app/rtc

FROM scratch AS prod-stage

COPY --from=build-stage /app/rtc /app/rtc 

# TLS certificate files.
COPY --from=build-stage /etc/ssl/certs/ /etc/ssl/certs/

USER 1000:1000

EXPOSE 8080

CMD ["/app/rtc"]
