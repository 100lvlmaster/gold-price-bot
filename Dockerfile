# syntax=docker/dockerfile:1

FROM golang:1.25 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /out/gold-price-api ./cmd/api/main.go

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/gold-price-api /gold-price-api
EXPOSE 8080
ENTRYPOINT ["/gold-price-api"]
