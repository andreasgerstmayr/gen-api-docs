FROM golang:1.21 AS build
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 go build -o /crd-to-cr


FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /crd-to-cr .
USER nonroot:nonroot
ENTRYPOINT ["/crd-to-cr"]
