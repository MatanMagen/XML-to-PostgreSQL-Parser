FROM golang:1.22 AS build

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main /go/src/app/


FROM gcr.io/distroless/static-debian12
# FROM alpine:latest # use it for development purposes to provide shell access

COPY --from=build /go/src/app/main /

CMD ["/main"]