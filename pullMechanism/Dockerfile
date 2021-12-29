FROM golang:1.16-buster AS build

FROM build AS dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . .

RUN go mod download

RUN go build -o pullmechanism
CMD ["./pullmechanism"]
