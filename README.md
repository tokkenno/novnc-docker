# novnc-docker
Docker image with NoVNC

- Native lightweight HTTP server.
- Native tcp-to-websocket proxy in replace of [WebSockify](https://github.com/novnc/websockify).
- WebUI to handle and manage multiple VNC connections.
- All in less than 6Mb!!!

## Configuration

The configuration file is located in `/etc/novnc/manager.json` and has the stucture:

```json
{
  "port": 8084,
  "servers": [
    {
      "name": "Localhost",
      "host": "localhost",
      "port": 5900
    }
  ]
}
```

- **Port:** Define the listen port for HTTP entry point where the NoVNC manager is served.
- **Servers:** Define zero, one or more remote VNC servers. The manager will create tcp-to-websocket proxies for all and will show a quick-connect button in manager interface for each of them.