FROM golang:1.21.7 AS builder 
    
LABEL maintainer "Marco Ferretti <mferretti93@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Starting a new stage from scratch 

FROM alpine:3.14

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

WORKDIR /root/

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]
