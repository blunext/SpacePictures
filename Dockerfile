FROM golang:1.16 as goBuilder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -ldflags "-s -w" -v -o server

FROM alpine:3
COPY --from=goBuilder /app/server /server
ENV API_KEY DEMO_KEY
ENV CONCURRENT_REQUESTS 5
ENV PORT 8080

EXPOSE $PORT
CMD ["/server"]