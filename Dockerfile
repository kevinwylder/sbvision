FROM golang:1.12

RUN go get -u github.com/aws/aws-sdk-go/...
RUN go get -u github.com/go-sql-driver/mysql

RUN apt-get update && apt-get install -y \
    curl

RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl && \
    chmod a+rx /usr/local/bin/youtube-dl

CMD go run github.com/kevinwylder/sbvision/cmd/server