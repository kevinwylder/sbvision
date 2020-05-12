FROM golang:1.12

RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/jmespath/go-jmespath
RUN go get github.com/gorilla/websocket
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/lestrrat-go/jwx/jwk

RUN apt-get update && apt-get install -y \
    curl ffmpeg

RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl && \
    chmod a+rx /usr/local/bin/youtube-dl

RUN go get github.com/kevinwylder/sbvision

ENTRYPOINT [ "go", "run", "github.com/kevinwylder/sbvision/cmd/server" ]

