FROM alpine:3.12 as server

RUN apk add --no-cache go

WORKDIR /build
COPY src /build
RUN CGO_ENABLED=0 GOOS=$(echo $TARGETPLATFORM| cut -d'/' -f 1) GOARCH=$(echo $TARGETPLATFORM| cut -d'/' -f 2) go build -a -installsuffix cgo -ldflags="-w -s" -o /build/novnc-manager

FROM alpine:3.12 as client

RUN apk add --no-cache git nodejs npm

RUN npm install -g webpack-cli @angular/cli

WORKDIR /build/client
COPY src/client /build/client
RUN npm install
RUN ng build --prod

WORKDIR /build
RUN git clone --branch v1.2.0 https://github.com/novnc/noVNC.git novnc
WORKDIR /build/novnc
RUN npm install && ./utils/use_require.js --with-app --as commonjs

FROM alpine:3.12 as runtime

RUN mkdir -p /app/novnc && mkdir -p /app/client
COPY --from=server /build/novnc-manager /app/novnc-manager
COPY --from=client /build/client/dist/client /app/client/
COPY --from=client /build/novnc/build /app/novnc/

WORKDIR /app
EXPOSE 8084
ENTRYPOINT [ "/app/novnc-manager" ]
