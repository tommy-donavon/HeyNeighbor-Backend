import mongoose from 'mongoose';

export interface IChat extends mongoose.Document {
  room_name: string;
  property_code: string;
  messages: Array<IMessage>;
}
export interface IMessage extends mongoose.Document {
  user_name: string;
  image_uri: string;
  message: string;
  time: Date
}


const MessageSchema: mongoose.Schema = new mongoose.Schema({
  user_name: {
    type: String,
    required: true,
  },
  image_uri: {
    type: String,
  },
  message: {
    type: String,
  },
  time: {
    type: Date,
    default: Date.now()
  }
});

const ChatSchema: mongoose.Schema = new mongoose.Schema({
  room_name: {
    type: String,
    required: true,
  },
  property_code: {
    type: String,
    required: true,
  },
  messages: [
    {
      type: mongoose.SchemaTypes.ObjectId,
      ref: 'Messages',
    },
  ],
});

ChatSchema.index({property_code: 1, room_name:1},{unique:true})
export const Message: mongoose.Model<IMessage> = mongoose.model(
  'Messages',
  MessageSchema,
);
export const Chat: mongoose.Model<IChat> = mongoose.model('Chats', ChatSchema);

export const newRooms =  (propCode: string, roomNames: Array<string>): IChat[] => {
  let output:IChat[] = []
  roomNames.forEach(async (rn) => {
     Chat.create({
     room_name: rn,
     property_code: propCode,
   }, (err, chat) => {
     if(err) console.error(err)
     output.push(chat)
   });
 });
 return output
};