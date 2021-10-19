import { Server, Socket } from 'socket.io';
import http from 'http';
import { getProperty, getUser } from '../utils/request.js';
import { IProperty, Chat, IChat, newRooms, Message } from '../models/index.js';
import url from 'url';

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

    const ns = io.of(/\/\w+/);
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
        if (!validChannel) {
          console.log('i disconnected');
          socket.disconnect(true);
        }
        let pastChats: IChat[] = (await Chat.find({
          property_code: prop.server_code,
        })) as IChat[];
        if (pastChats.length === 0) {
          pastChats = newRooms(prop.server_code, prop.channels);
        }
        socket.on('msg', async (msg) => {
          socket.broadcast.to(room as string).emit('msg', msg);
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
                  console.log(c.messages);
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
    return io;
  };
}
