import {Component, OnInit} from '@angular/core';
import {HttpClient} from "@angular/common/http";

interface ServerConfig {
  name: string;
  host: string;
  port: number;
  proxy: string;
}

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  public servers: ServerConfig[] = [];
  public loadError: string = null;

  constructor(
    private readonly http: HttpClient
  ) {
  }

  ngOnInit(): void {
    this.http.get<ServerConfig[]>('/servers').subscribe((servers) => {
      this.servers = servers;
      this.servers.push({name: 'Brave 2', host: 'localhost', port: 6756, proxy: '/ws/d'});
      this.servers.push({name: 'Brave 3', host: 'localhost', port: 6756, proxy: '/ws/d'});
      this.servers.push({name: 'Brave 4', host: 'localhost', port: 6756, proxy: '/ws/d'});
      this.servers.push({name: 'Brave 5', host: 'localhost', port: 6756, proxy: '/ws/d'});
      this.loadError = null;
    }, (error) => {
      console.error('Load error:', error);
      this.loadError = 'Unknown load error';
    });
  }

  getUrl(server: ServerConfig) {
    const params: string[] = [];
    if (server.proxy != null && server.proxy.length > 0) {
      if (server.proxy.startsWith('/')) {
        server.proxy = server.proxy.substring(1);
      }
      params.push(`path=${server.proxy}`);
    }
    params.push(`autoconnect=true`);
    return '/vnc/vnc.html?' + params.join('&');
  }
}
