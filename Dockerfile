FROM golang:1.17 as build

WORKDIR /go/src/

COPY . .

RUN go mod download

RUN make build-server && make build-client 

FROM gcr.io/distroless/base-debian11

COPY --from=build /go/src/wsp_server /
COPY --from=build /go/src/wsp_client / 
