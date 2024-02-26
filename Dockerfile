FROM golang:1.21-bookworm AS builder
WORKDIR /build
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o hcloud-controller ./cmd/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /build/hcloud-controller /
CMD ["/hcloud-controller"]
