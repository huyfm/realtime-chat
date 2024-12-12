FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate

RUN go build -o /app/rtc -tags netgo ./cmd

RUN chown 1000:1000 /app/rtc

FROM alpine:3.21 AS prod-stage

COPY --from=build-stage /app/rtc /app/rtc 

USER 1000:1000

EXPOSE 8080

CMD ["/app/rtc"]
