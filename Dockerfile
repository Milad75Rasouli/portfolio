# sudo docker build -t portfolio:latest .

FROM golang:1.22.1-alpine3.19 as build

WORKDIR /app
RUN apk add --update build-base
COPY ./go.mod .
COPY ./go.sum .
RUN export GOPROXY=https://gocenter.io,direct
RUN go mod download
COPY . /app
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildvcs=false -trimpath -ldflags="-w -s" -o ./bin/portfolio ./main.go

FROM alpine:3.19.1
WORKDIR /app
COPY --from=build /app/bin/portfolio /app
COPY --from=build /app/frontend /app/frontend
EXPOSE 5000
ENTRYPOINT ["/app/portfolio"]
