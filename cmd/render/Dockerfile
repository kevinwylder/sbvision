FROM golang:1.12

RUN apt-get update && apt-get install -y \
    ffmpeg libgl1-mesa-dev xvfb xorg-dev

RUN go get golang.org/x/image/bmp
RUN go get github.com/go-gl/gl/v3.3-core/gl
RUN go get -u github.com/go-gl/glfw/v3.3/glfw
RUN go get github.com/kevinwylder/sbvision
RUN go install github.com/kevinwylder/sbvision/cmd/render
