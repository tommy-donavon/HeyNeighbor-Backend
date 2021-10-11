import { Server, Socket } from 'socket.io';
import http from "http"
import { getProperty, getUser } from '../utils/request';
import {IProperty, Chat, IChat, newRooms} from '../models'
import url from 'url';

export class SocketServer {
  static SetupServer = (server: http.Server): Server => {
    let io = new Server({
      path: '/',
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
          if (!validChannel) socket.disconnect(true);
          let pastChats:IChat = await Chat.findOne({server_code: prop.server_code}) as IChat
          if(!pastChats) {
            pastChats = await Chat.create({
              server_code: prop.server_code,
              rooms: newRooms
            })
          }
          socket.join(room as string);
          console.log(
            `socket joined room ${room} in namespace ${socket.nsp.name.substr(
              1,
            )}`,
          );
          socket.on('msg', async (msg) => {
            console.log(msg);
            socket.broadcast.to(room as string).emit('msg', msg);
            pastChats.rooms.forEach(c => {
              if(c.name === room as string) {
                c.chats.push({
                  user_name: user.username,
                  image_uri: msg.image_uri,
                  message: msg.message
                })
              }
            })
            await Chat.findOneAndUpdate({server_code: pastChats.server_code}, pastChats)
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
      return io
  }
}
