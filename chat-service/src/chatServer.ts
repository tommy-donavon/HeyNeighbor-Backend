import express, { Application, Request, Response } from 'express';
import http from 'http';
import { Server } from 'socket.io';
import redisAdapter from '@socket.io/redis-adapter';
import redis from 'redis';
import {
  newConsulDetails,
  registerService,
  deregisterService,
} from './register/register.js';
import { SocketServer } from './handlers/socketServer.js';
import router from './handlers/http.js';

export class ChatServer {
  public static readonly PORT: number = 8080;
  private app: Application;
  private server: http.Server;
  private io: Server;
  private port: string | number;

  constructor() {
    this.app = express();
    this.port = process.env.PORT || ChatServer.PORT;
    this.app.use((req: Request, res: Response, next) => {
      res.setHeader('Content-Type', 'application/json');
      next();
    });
    this.app.use((req: Request, res: Response, next) => {
      res.header('Access-Control-Allow-Origin', '*');
      res.header(
        'Access-Control-Allow-Headers',
        'Origin, X-Requested-With, Content-Type, Accept',
      );
      next();
    });
    this.app.use(router);
    this.server = http.createServer(this.app);
    this.io = SocketServer.SetupServer(this.server);

    const pubClient = redis.createClient({
      host: process.env.REDIS_HOST,
      port: Number(process.env.REDIS_PORT),
    });
    const subClient = pubClient.duplicate();
    this.io.adapter(redisAdapter.createAdapter(pubClient, subClient));

    this.server.listen(this.port, () => {
      var details = newConsulDetails('chat-service', <number>this.port);
      registerService(details);
      console.log(`server is running on port ${this.port}`);

      process.on('SIGTERM', () => {
        deregisterService(<string>details.id);
        process.exitCode = 1;
      });
    });
  }

  public getApp(): Application {
    return this.app;
  }
}
