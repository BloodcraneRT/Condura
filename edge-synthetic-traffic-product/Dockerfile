FROM node:20 AS ui-builder

WORKDIR /app

# Install pnpm per project memory requirements
RUN npm install -g pnpm

# Copy UI package files for layer caching
COPY ui/package.json ui/package-lock.json* ui/pnpm-lock.yaml* ./ui/

# Install dependencies
WORKDIR /app/ui
RUN pnpm install

# Copy the rest of the UI source
COPY ui/ ./

# Build UI
RUN pnpm build

# Go Builder
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Install make
RUN apk add --no-cache make git

# Copy everything
COPY . .

# Remove previously copied UI and use the one built from UI builder
RUN rm -rf ui server/ui-dist

# Copy the built UI
COPY --from=ui-builder /app/ui/dist /app/server/ui-dist

# Build the server binary
RUN cd server && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/bin/edgesynth

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required system tools (tcpreplay, traceroute/tracert equivalent) as per memory
RUN apk add --no-cache tcpreplay traceroute ca-certificates tzdata

COPY --from=go-builder /app/bin/edgesynth /app/edgesynth

EXPOSE 8080

ENTRYPOINT ["/app/edgesynth"]
