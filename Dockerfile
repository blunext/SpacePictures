FROM golang:1.16 as goBuilder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -ldflags "-s -w" -v -o server

FROM alpine:3
COPY --from=goBuilder /app/server /server
EXPOSE 8080
CMD ["/server"]