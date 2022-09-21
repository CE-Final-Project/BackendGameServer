FROM golang:alpine3.16 as build

WORKDIR /app

COPY ../ ./

RUN cd account && go mod download


RUN cd account && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o account_service ./cmd/main.go

FROM scratch

COPY --from=build /app/account/account_service /account_service
COPY ../account/config /config
EXPOSE 5003

CMD ["./account_service"]