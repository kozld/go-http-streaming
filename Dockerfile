FROM golang:1.19.3-alpine AS builder
RUN apk update && apk add --no-cache git make
WORKDIR /app
COPY . .
RUN make build

FROM scratch
COPY --from=builder /app/build/app app/
WORKDIR app
ENTRYPOINT ["./app"]
