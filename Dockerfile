# Etapa 1: Build
FROM golang:1.24-alpine AS build

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./core/infra/config/credentials.json ./core/infra/config/credentials.json


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /queue ./cmd/main.go

# Etapa 2: Imagem final
FROM scratch

COPY --from=build /queue /queue
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/core/infra/config/credentials.json /core/infra/config/credentials.json
COPY .env /

EXPOSE 3003

CMD ["./queue"]

#docker run --rm -p 3003:3003 --name ec-queue --env-file .env ec-queue comando pra rodar o container