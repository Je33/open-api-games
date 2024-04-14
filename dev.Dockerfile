FROM golang:1.22-alpine AS dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY . ./
RUN go mod download

CMD ["air", "-c", "api.air.toml"]
