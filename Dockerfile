FROM golang:latest

WORKDIR /app
ADD . /app/

EXPOSE 8080

CMD ["go", "build", "yt_fetch.go"]
