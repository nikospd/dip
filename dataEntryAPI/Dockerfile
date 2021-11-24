FROM golang:1.16-buster AS build

FROM build AS dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . .

RUN go mod download


EXPOSE 1324
RUN go build -o server
CMD ["./server"]
