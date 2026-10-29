FROM  golang:1.27 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/httpserver

FROM scratch

COPY --from=builder /bin/server /server

EXPOSE 8080

ENTRYPOINT [ "/server" ]
