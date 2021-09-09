import express, { Application, Request, Response } from "express";
import http from "http";
import { Namespace, Server, Socket } from "socket.io";
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
    this.io = new Server().listen(this.server);

    this.server.listen(this.port, () => {
      var details = newConsulDetails("chat-service", <number>this.port);
      registerService(details);
      console.log(`server is running on port ${this.port}`);

      process.on("SIGTERM", () => {
        deregisterService(<string>details.id);
        process.exitCode = 1;
      });
    });

    this.io.sockets.on("connection", (socket: Socket) => {
      var handShackURL = url.parse(socket.handshake.url, true);
      var ns = handShackURL.query.ns;
      var room = handShackURL.query.room;

      if (typeof room !== "string" || typeof ns !== "string") {
        socket.disconnect(true);
      }
      var validNs = props.filter(
        (p) =>
          p.serverCode === (ns as string) && p.Channels.includes(room as string)
      );
      if (validNs.length !== 1) {
        socket.disconnect(true);
        return;
      }

      this.io.of(ns as string).on("connection", (s2: Socket) => {
        s2.join(room as string);
        console.log(`socket joined room: ${room} in namespace: ${ns}`);

        s2.on("disconnect", (reason: string) => {
          s2.leave(room as string);
          console.log(reason);
        });
      });
      socket.on("disconnect", (reason: string) => console.log(reason));
      socket.on("error", (err: Error) => {
        console.log(err);
        socket.disconnect(true);
      });
    });
  }

  public getApp(): Application {
    return this.app;
  }
}
