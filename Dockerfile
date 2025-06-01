FROM golang:1.24

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build only the main package
RUN go build -v -o /usr/local/bin/app ./cmd

CMD ["app"]
