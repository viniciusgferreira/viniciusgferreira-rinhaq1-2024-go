FROM golang:1.21 AS BUILD

WORKDIR /rinha

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

FROM build AS final

COPY . .
RUN go build -o app

CMD ["./app"]