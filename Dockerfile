FROM golang:1.15 as server

WORKDIR /build
COPY src /build
RUN CGO_ENABLED=0 GOOS=$(echo $TARGETPLATFORM| cut -d'/' -f 1) GOARCH=$(echo $TARGETPLATFORM| cut -d'/' -f 2) go build -a -installsuffix cgo -ldflags="-w -s" -o /build/novnc-manager

FROM node:12 as client

RUN npm install -g webpack-cli @angular/cli

WORKDIR /build/client
COPY src/client /build/client
RUN npm install
RUN ng build --prod

WORKDIR /build
RUN git clone --branch v1.2.0 https://github.com/novnc/noVNC.git novnc
WORKDIR /build/novnc
RUN npm install && ./utils/use_require.js --with-app --as commonjs

FROM alpine:latest as runtime

RUN mkdir -p /app/novnc && mkdir -p /app/client
COPY --from=server /build/novnc-manager /app/novnc-manager
COPY --from=client /build/client/dist/* /app/client/
COPY --from=client /build/novnc/* /app/novnc/

WORKDIR /app
EXPOSE 8084
ENTRYPOINT [ "/app/novnc-manager" ]
