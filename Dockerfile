FROM golang:1.12

RUN go get -u github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/websocket
RUN go get github.com/codegangsta/gin
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/lestrrat-go/jwx/jwk
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/jmespath/go-jmespath

RUN apt-get update && apt-get install -y \
    curl ffmpeg

RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl && \
    chmod a+rx /usr/local/bin/youtube-dl
