FROM golang:1.15 as server

WORKDIR /build
COPY src /build
RUN CGO_ENABLED=0 GOOS=$(echo $TARGETPLATFORM| cut -d'/' -f 1) GOARCH=$(echo $TARGETPLATFORM| cut -d'/' -f 2) go build -a -installsuffix cgo -ldflags="-w -s" -o /build/novnc-manager

FROM node:14 as client

RUN apk --no-cache add git

WORKDIR /build/client
COPY src/client /build/client
RUN npm install
RUN webpack build --production

WORKDIR /build/novnc
RUN git clone --branch v1.2.0 https://github.com/novnc/noVNC.git .
RUN npm install && ./utils/use_require.js --with-app --as commonjs

FROM alpine:latest as runtime

WORKDIR /app
RUN mkdir novnc client
COPY --from=server /build/novnc-manager /app/novnc-manager
COPY --from=client /build/client/dist/* /app/client
COPY --from=client /build/novnc/* /app/novnc

EXPOSE 8084
ENTRYPOINT [ "/app/novnc-manager" ]