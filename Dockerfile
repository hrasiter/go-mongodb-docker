FROM golang AS builder
WORKDIR /build/api
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o api

FROM alpine
WORKDIR /root
COPY --from=builder /build/api/api .
COPY --from=builder /build/api/environment.env .
EXPOSE 8080
CMD [ "./api" ]
 