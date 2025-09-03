FROM golang:1.22-alpine

# Install dependencies
RUN apk add --no-cache git curl

# Install air
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

WORKDIR /app

# Disable VCS stamping
ENV GOFLAGS="-buildvcs=false"

# Tambah user (opsional)
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8000

CMD ["air"]
