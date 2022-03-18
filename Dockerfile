FROM golang:alpine

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go mod download
CMD [ "go", "test", "-v" ]