FROM golang:1.17
WORKDIR /app
ENV HTTP_PORT=8080
EXPOSE ${HTTP_PORT}
COPY go.mod ./
COPY *.go ./
RUN go build -o /http-server
CMD ["/http-server"]