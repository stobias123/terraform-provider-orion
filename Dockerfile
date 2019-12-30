FROM golang:1.13-buster as builder

WORKDIR /go/terraform-solarwinds
#COPY go.mod .
#COPY go.sum .
#RUN go mod download

COPY . .
RUN go build -o /out/terraform-provider-orion

FROM debian:buster

COPY --from=builder /out/terraform-provider-orion /