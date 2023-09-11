FROM golang:1.21.0

WORKDIR /pavosql

COPY ./ ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/pavosql

EXPOSE 1758

CMD [ "./pavosql" ]
