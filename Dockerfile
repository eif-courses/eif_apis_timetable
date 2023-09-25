FROM golang:1.17

WORKDIR /go/src/eif_apis_timetable
COPY . .
RUN go get -d -v ./...
RUN go build

CMD ["/go/src/eif_apis_timetable/eif_apis_timetable"]
