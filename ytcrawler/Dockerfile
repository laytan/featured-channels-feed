# Start with a base of golang on alpine linux
FROM golang:1.15.6-alpine3.12

# Allow specifying the frontend url
ARG FRONTEND_URL
ENV FRONTEND_URL=$FRONTEND_URL

# Make a /app directory to put our go application
RUN mkdir /app

# Run all commands after this one from /app
WORKDIR /app

# Copy our code
COPY . .

# Download dependencies
RUN go mod download

# Install Chromium and it's dependencies
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" > /etc/apk/repositories \
    && echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
    && echo "http://dl-cdn.alpinelinux.org/alpine/v3.12/main" >> /etc/apk/repositories \
    && apk upgrade -U -a \
    && apk add \
    libstdc++ \
    chromium \
    harfbuzz \
    nss \
    freetype \
    ttf-freefont \
    font-noto-emoji \
    wqy-zenhei \
    && rm -rf /var/cache/* \
    && mkdir /var/cache/apk

# Set env variables for chromium
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

# Build the application
WORKDIR /app/cmd/ytcrawler
RUN go build

# Start the app
CMD ["./ytcrawler"]
