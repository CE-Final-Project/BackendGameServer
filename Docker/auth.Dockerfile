FROM golang:alpine3.16 as build

WORKDIR /app

COPY ../ ./

RUN cd authentication && go mod download


RUN cd authentication && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth_service ./cmd/main.go

FROM scratch

COPY --from=build /app/authentication/auth_service /auth_service
COPY ../authentication/config /config
EXPOSE 5003

CMD ["./auth_service"]