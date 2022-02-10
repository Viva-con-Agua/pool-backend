FROM golang


WORKDIR /app

ADD . /app

RUN go install

CMD ["go", "run", "server.go"]
