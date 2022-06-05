FROM golang:1.18-alpine
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build main.go

RUN mv main /usr/local/bin/entrypoint && chmod +x /usr/local/bin/entrypoint

ENTRYPOINT ["entrypoint"]