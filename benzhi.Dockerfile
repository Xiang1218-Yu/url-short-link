FROM golang:1.26.5-alpine

ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

WORKDIR /workspace
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build ./...

ENTRYPOINT ["go", "run", "./cmd/shortlink"]
CMD ["--input", "examples/links.json", "--code", "docs", "--at", "2026-08-15T12:00:00Z"]
