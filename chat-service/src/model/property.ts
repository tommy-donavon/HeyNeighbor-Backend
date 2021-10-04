import mongoose, {Schema, Document} from "mongoose";

export interface IProperty{
  server_code: string;
  channels: string[];
}


