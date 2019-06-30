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
RUN mkdir /squadxml

# default config
ENV SQUADXML_PATH=/squadxml
ENV SQUADXML_HOST=0.0.0.0:80
ENV SQUADXML_DB_USER=user
ENV SQUADXML_DB_PASSWORD=password
ENV SQUADXML_DB_HOST=host
ENV SQUADXML_DB=db

VOLUME ["/squadxml"]
CMD ["/app/main"]
