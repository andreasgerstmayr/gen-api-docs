FROM golang:1.21 AS build
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /gen-api-docs


FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /gen-api-docs .
USER nonroot:nonroot
ENTRYPOINT ["/gen-api-docs"]
