import { Schema, model, Model, Document } from "mongoose";

export interface IChat extends Document {
    server_code: string;
    rooms:Array<IRoom>
}

interface IRoom {
    name:string;
    chats:Array<IMessage>
}
interface IMessage {
    user_name:string;
    image_uri:string;
    message:string
}
export const newRooms = (roomNames:Array<string>):Array<IRoom> => {
    let output:Array<IRoom> = []
    roomNames.forEach(rn => {
        var room:IRoom = {
            name: rn,
            chats: [],
        }
        output.push(room)
    })
    return output
}

const ChatSchema:Schema = new Schema({
    property_code: {
        type:String,
        required:true,
        unique:true
    },
    rooms:[
        {
            room_name: {
                type:String,
                required:true
            },
            chats:[
                {
                    user_name:{
                        type:String,
                        required:true
                    },
                    image_uri:String,
                    message:String
                }
            ]
        }
    ]
})
export const Chat: Model<IChat> = model('Chats', ChatSchema)