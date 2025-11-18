FROM --platform=linux/amd64 golang:1.25.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o example _cmd/main.go

FROM --platform=linux/amd64 alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Bangkok
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


WORKDIR /app
COPY --from=builder /app/example .
EXPOSE 8080
CMD ["./example"]