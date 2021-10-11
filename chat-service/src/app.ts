import { ChatServer } from "./chatServer.js";
import Connect from "./dbconnect.js";


Connect(process.env.MONGO_URI as string)
const app = new ChatServer().getApp()

export {app};