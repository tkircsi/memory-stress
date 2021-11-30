FROM golang:1.16 as build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/mstress .

FROM scratch as bin
COPY --from=build /app/mstress /app/mstress

ENV LARGE_FILE="/tmp/large.file"
EXPOSE 8080
ENTRYPOINT [ "/app/mstress" ]