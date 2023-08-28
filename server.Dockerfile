FROM golang:latest as builder

WORKDIR app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o server


#FROM alpine
#WORKDIR /
#COPY --from=builder /app/server .

CMD ["./server"]