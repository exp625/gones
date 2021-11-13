FROM golang:1.16 AS build

RUN apt-get update
RUN apt-get install -y libgl1-mesa-dev xorg-dev gcc-multilib gcc-mingw-w64 libasound2-dev
WORKDIR /usr/src/nes
RUN mkdir dist
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify
VOLUME /root/.cache/go-build
CMD CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
go build -o ./bin/nes.exe ./main.go