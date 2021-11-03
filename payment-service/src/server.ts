import express, { Application } from 'express';
import http, {Server} from 'http';

export default class AppServer {
  public static readonly PORT: number = 8080;
  private app: Application;
  private server: Server;
  private port: string | number;

  constructor() {
    this.app = express();
    this.port = process.env.PORT || AppServer.PORT;
    this.server = http.createServer(this.app);
    this.server.listen(this.port, () => {
      console.log(`Server is listening on port: ${this.port}`);
    });
  }

  public getApp(): Application {
    return this.app;
  }
}
