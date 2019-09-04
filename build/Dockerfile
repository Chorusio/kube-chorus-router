FROM golang:alpine AS builder
WORKDIR $GOPATH/src/k8s-route-extender/
COPY cmd/ cmd
COPY version/ version
COPY vendor/ vendor  
RUN go build -o /go/bin/kube-chorus-router ./cmd/main.go

FROM alpine
COPY --from=builder /go/bin/kube-chorus-router /go/bin/kube-chorus-router
ENTRYPOINT ["/go/bin/kube-chorus-router"]
