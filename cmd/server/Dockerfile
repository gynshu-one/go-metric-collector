FROM arm64v8/golang:latest

WORKDIR /app

COPY ../../go.mod ./
COPY ../../go.sum ./

RUN go mod download && go mod verify
COPY ../.. .
RUN go test ./...
#RUN go build -o  ./cmd/server/server ./cmd/server/server.go

