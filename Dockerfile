FROM golang:1.24

WORKDIR /antibruteforce

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN ls ./

RUN go build -o /antibuteforce/build/service /antibruteforce/cmd/service

EXPOSE 8080

ENTRYPOINT ["/antibuteforce/build/service", "--conf=config/conf.yaml"]