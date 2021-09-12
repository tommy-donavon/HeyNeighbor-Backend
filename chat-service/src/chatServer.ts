import express, { Application } from "express";
import http from "http";
import { Server, Socket } from "socket.io";
import redisAdapter from "@socket.io/redis-adapter";
import redis from "redis";
import url from "url";
import { IProperty } from "./model";
import {
  newConsulDetails,
  registerService,
  deregisterService,
} from "./register/register.js";

export class ChatServer {
  public static readonly PORT: number = 8080;
  private app: Application;
  private server: http.Server;
  private io: Server;
  private port: string | number;

  constructor(props: IProperty[]) {
    this.app = express();
    this.port = process.env.PORT || ChatServer.PORT;
    this.server = http.createServer(this.app);
    this.io = new Server({ path: "/" }).listen(this.server);

    const pubClient = redis.createClient({
      host: process.env.REDIS_HOST,
      port: Number(process.env.REDIS_PORT),
    });
    const subClient = pubClient.duplicate();
    this.io.adapter(redisAdapter.createAdapter(pubClient, subClient));

    this.server.listen(this.port, () => {
      var details = newConsulDetails("chat-service", <number>this.port);
      registerService(details);
      console.log(`server is running on port ${this.port}`);

      process.on("SIGTERM", () => {
        deregisterService(<string>details.id);
        process.exitCode = 1;
      });
    });

    const ns = this.io.of(/\/\w+/);
    ns.on("connection", (socket: Socket) => {
      var room = url.parse(socket.handshake.url, true).query.room;

      if (typeof room !== "string") {
        socket.disconnect(true);
      }
      var validNs = props.filter(
        (p) =>
          p.serverCode === socket.nsp.name.substr(1) &&
          p.Channels.includes(room as string)
      );
      if (validNs.length !== 1) {
        socket.disconnect(true);
        return;
      }
      socket.join(room as string);
      console.log(
        `socket joined room ${room} in namespace ${socket.nsp.name.substr(1)}`
      );

      socket.on("disconnect", (reason: string) => console.log(reason));
      socket.on("error", (err: Error) => {
        console.log(err);
        socket.leave(room as string);
        socket.disconnect(true);
      });
    });
  }

  public getApp(): Application {
    return this.app;
  }
}
