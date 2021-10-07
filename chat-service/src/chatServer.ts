import express, { Application } from 'express';
import http from 'http';
import { Server, Socket } from 'socket.io';
import redisAdapter from '@socket.io/redis-adapter';
import redis from 'redis';
import url from 'url';
import { getUser, getProperty } from './utils/request.js';
import {
  newConsulDetails,
  registerService,
  deregisterService,
} from './register/register.js';
import { IProperty } from './model/property.js';

export class ChatServer {
  public static readonly PORT: number = 8080;
  private app: Application;
  private server: http.Server;
  private io: Server;
  private port: string | number;

  constructor() {
    this.app = express();
    this.port = process.env.PORT || ChatServer.PORT;
    this.server = http.createServer(this.app);
    this.io = new Server({
      path: '/',
      transports: ['websocket', 'polling'],
      cors: {
        origin: '*',
        methods: ['GET', 'POST'],
        credentials: true,
      },
      allowEIO3: true,
    }).listen(this.server);

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

    const ns = this.io.of(/\/\w+/);
    ns.on('connection', async (socket: Socket) => {
      try {
        var room = url.parse(socket.handshake.url, true).query.room;
        var server_code = socket.nsp.name.replace(/[/]/g, '');
        const user = await getUser(socket.request.headers.authorization);

        const prop = await getProperty(
          server_code,
          socket.request.headers.authorization as string,
        );
        var validChannel = (prop as IProperty).channels.includes(
          room as string,
        );
        if (!validChannel) socket.disconnect(true);
        socket.join(room as string);
        console.log(
          `socket joined room ${room} in namespace ${socket.nsp.name.substr(
            1,
          )}`,
        );
        socket.on('msg', (msg) => {
          console.log(msg);
          socket.broadcast.to(room as string).emit('msg', msg);
        });

        socket.on('disconnect', (reason: string) => console.log(reason));
        socket.on('error', (err: Error) => {
          console.error(err);
          socket.leave(room as string);
          socket.disconnect(true);
        });
      } catch (err) {
        console.error(err);
        socket.disconnect(true);
      }
    });
  }

  public getApp(): Application {
    return this.app;
  }
}
