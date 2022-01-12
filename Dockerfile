FROM golang:alpine as builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build

FROM golang:alpine as app
WORKDIR /app
COPY --from=builder /build/nextcloud-exporter ./
COPY --from=builder /build/templates ./templates
CMD [ "./nextcloud-exporter" ]