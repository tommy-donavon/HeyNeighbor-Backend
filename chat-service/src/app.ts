import { ChatServer } from "./chatServer.js";
import { IProperty } from "./model/index.js";
import Connect from "./dbconnect.js";


//TODO get connection string from env.
// Connect("mongodb://localhost:27017/chat?ssl=false")
let props: IProperty[] = [
    {
        serverCode: "test",
        Channels:[
            "general",
            "announcements",
            "events"
        ]
    }
]
let app = new ChatServer().getApp()

export {app};