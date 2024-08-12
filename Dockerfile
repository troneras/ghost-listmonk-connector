# Build the Next.js app
FROM node:20.10.0 AS nextjs-builder
WORKDIR /app
COPY ui/ ./
RUN npm ci
RUN npm run build

# Build the Go app
FROM golang:1.21.6 AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates nodejs npm && rm -rf /var/lib/apt/lists/*

WORKDIR /root/
COPY --from=go-builder /app/main .
COPY --from=nextjs-builder /app .
COPY start.sh .

ENV GIN_MODE=release
ENV NODE_ENV=production

EXPOSE 8808 3000

CMD ["./start.sh"]