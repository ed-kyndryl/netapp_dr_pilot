FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/quotes-web .

FROM gcr.io/distroless/static:nonroot
USER 65532:65532
COPY --from=build /out/quotes-web /quotes-web

EXPOSE 8080
ENTRYPOINT ["/quotes-web"]
