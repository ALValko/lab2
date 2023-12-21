FROM golang:1.21.0-alpine AS build
WORKDIR /build
COPY . .
RUN go build -v -o ./main .

FROM alpine:latest
WORKDIR /app
COPY --from=build /build ./
CMD ["./main"]