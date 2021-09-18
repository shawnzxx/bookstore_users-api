## Multi-stage builds
FROM golang:1.17-alpine AS builder
WORKDIR /app
#here use absoluteDir pattern
COPY . /app/

RUN go mod tidy
RUN go mod verify

# CGO has to be disabled for FROM scratch: CGO_ENABLED=0
# https://stackoverflow.com/questions/52640304/standard-init-linux-go190-exec-user-process-caused-no-such-file-or-directory
# https://stackoverflow.com/questions/62817082/how-does-cgo-enabled-affect-dynamic-vs-static-linking
# https://www.geeksforgeeks.org/static-and-dynamic-linking-in-operating-systems/
#here use relativeDir pattern
RUN CGO_ENABLED=0 go build -o bookstore-user-api

## Deploy and run binary
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bookstore-user-api .
EXPOSE 8081
ENTRYPOINT ["./bookstore-user-api"]