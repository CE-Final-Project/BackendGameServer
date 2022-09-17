FROM golang:alpine3.16 as build

WORKDIR /app

COPY ../ ./

RUN cd gateway && go mod download


RUN cd gateway && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway_service ./cmd/main.go

FROM scratch

COPY --from=build /app/gateway/gateway_service /gateway_service
COPY ../gateway/config /config
EXPOSE 8001

CMD ["./gateway_service"]