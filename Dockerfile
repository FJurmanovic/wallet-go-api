FROM golang:alpine as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/api ./cmd/api/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/bin/migrate /usr/bin/
COPY --from=builder /app/bin/api /usr/bin/

CMD ["migrate"]
ENTRYPOINT ["api"]
