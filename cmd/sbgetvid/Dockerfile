FROM golang:1.12

RUN apt-get update && \
    apt-get install -y ffmpeg

RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/jmespath/go-jmespath
RUN go get github.com/kevinwylder/sbvision

RUN go install github.com/kevinwylder/sbvision/cmd/sbgetvid


