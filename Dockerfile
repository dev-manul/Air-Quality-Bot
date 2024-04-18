FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/bot

ARG BOT_TOKEN
ARG AQICN_TOKEN

CMD ["/bot"]