FROM golang AS build
ADD . /go/src/github.com/Kugelschieber/squadxml
WORKDIR /go/src/github.com/Kugelschieber/squadxml
RUN apt update && \
	apt upgrade -y
ENV GOPATH=/go
RUN go build -ldflags "-s -w" main.go

FROM alpine
COPY --from=build /go/src/github.com/Kugelschieber/squadxml /app
WORKDIR /app

# default config
ENV SQUADXML_PATH=/squadxml
ENV SQUADXML_HOST=0.0.0.0:80

VOLUME ["/squadxml"]
CMD ["/app/main"]
