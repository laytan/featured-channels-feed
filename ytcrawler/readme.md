# ytcrawler

The server for featured-channels-feed.

Built with Golang and uses Docker for deployment.

## Dev

Make sure Docker and Docker Compose are installed on your system.

```bash
docker-compose up
```

This runs the server and rebuilds when files change.
It defaults to http://localhost:3000 as the front-end URL.
This can be changed in the docker-compose.yml file.

Run ```docker-compose build``` when you change a docker file.

## Build

Make sure Docker is installed on your system.

```bash
docker build --build-arg FRONTEND_URL=http://example.com -t featured-channels-crawler .
```

This builds an image called featured-channels-crawler with a fron-end url of http://example.com.

```bash
docker run -p 8080:8080 featured-channels-crawler
```

This runs the previously created image. And maps it to port 8080.