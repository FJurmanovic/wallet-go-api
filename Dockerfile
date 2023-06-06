FROM golang:alpine

WORKDIR /app 

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api ./cmd/api/main.go

EXPOSE 4000

CMD ["/api"]
