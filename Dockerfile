from golang as build

WORKDIR /tmp

COPY . .

RUN go build -o /bin/pavosql /tmp/cmd/pavosql

from debian:bookworm as run

USER 1000:1000
WORKDIR /

COPY --from=build /bin/pavosql pavosql

CMD [ "pavosql", "serve", "--port", "1234", "--file", "/data/pavosql/pavosql.db" ]
