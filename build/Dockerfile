FROM golang:alpine AS builder
WORKDIR $GOPATH/src/K8s-Route-Extender/
COPY cmd/ cmd
COPY version/ version
COPY vendor/ vendor  
RUN go build -o /go/bin/k8s-route-extender ./cmd/main.go

FROM alpine
COPY --from=builder /go/bin/k8s-route-extender /go/bin/k8s-route-extender
ENTRYPOINT ["/go/bin/k8s-route-extender"]
