import { Server, Socket } from 'socket.io';
import http from 'http';
import { getProperty, getUser } from '../utils/request.js';
import { IProperty, Chat, IChat, newRooms, Message } from '../models/index.js';
import url from 'url';
// import redisAdapter from '@socket.io/redis-adapter';
// import redis from 'redis';

export class SocketServer {
  static SetupServer = (server: http.Server): Server => {
    let io = new Server({
      path: '/socket/',
      transports: ['websocket', 'polling'],
      cors: {
        origin: '*',
        methods: ['GET', 'POST'],
        credentials: true,
      },
      allowEIO3: true,
    }).listen(server);
    // const pubClient = redis.createClient({
    //   host: process.env.REDIS_HOST,
    //   port: Number(process.env.REDIS_PORT),
    // });
    // const subClient = pubClient.duplicate();
    // io.adapter(redisAdapter.createAdapter(pubClient, subClient));

    const ns = io.of(/\/\w+/);
    ns.on('connection', async (socket: Socket) => {
      try {
        var room = url.parse(socket.handshake.url, true).query.room;
        var server_code = socket.nsp.name.replace(/[/]/g, '');

        const user = await getUser(socket.request.headers.authorization);
        console.log(`User: ${user.username} joined namespace: ${server_code} in room: ${room}`)

        const prop = await getProperty(
          server_code,
          socket.request.headers.authorization as string,
        );
        var validChannel = (prop as IProperty).channels.includes(
          room as string,
        );
        if (!validChannel) {
          let userRoom = (<string>room).split(':');
          if (
            (<IProperty>prop).tenants.filter(
              (t) => t.username.toUpperCase() === userRoom[0].toUpperCase() || t.username.toUpperCase() === userRoom[1].toUpperCase(),
            ).length !== 2
          ) {
            socket.disconnect(true);
            return;
          }
          Chat.create({
            room_name: <string>room,
            property_code: prop.server_code,
          });
        }
        let pastChats: IChat[] = (await Chat.find({
          property_code: prop.server_code,
        })) as IChat[];
        if (pastChats.length === 0) {
          pastChats = newRooms(prop.server_code, prop.channels);
        }
        socket.on('msg', async (msg) => {
          socket.to(room as string).emit('msg', msg);
          pastChats.forEach((c) => {
            if (c.room_name === <string>room) {
              Message.create(
                {
                  user_name: user.username,
                  image_uri: msg.image_uri,
                  message: msg.message,
                },
                async (err, im) => {
                  if (err) console.error(err);
                  if (!c.messages) c.messages = [];
                  c.messages.push(im);
                  Chat.updateOne(
                    { property_code: c.property_code, room_name: <string>room },
                    { $set: { messages: c.messages } },
                  ).catch((err) => console.error(err));
                },
              );
            }
          });
        });

        socket.on('disconnect', (reason: string) => console.log(reason));
      } catch (err) {
        console.error(err);
        socket.disconnect(true);
      }
    });
    return io;
  };
}
