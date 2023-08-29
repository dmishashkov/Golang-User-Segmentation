FROM golang:latest as builder

WORKDIR app

COPY . ./

RUN go mod download && go build -o server


#FROM alpine
#WORKDIR /
#COPY --from=builder /app/server .

CMD ["./server"]