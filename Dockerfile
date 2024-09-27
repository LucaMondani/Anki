
FROM golang:1.23-alpine3.19 AS build
WORKDIR /go/src/app
COPY . .
RUN go get .
RUN GOOS=linux go build -o ./bin/anki ./main.go

FROM alpine:3.19
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 1323
ENTRYPOINT /go/bin/anki 