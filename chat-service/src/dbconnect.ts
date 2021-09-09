import mongoose from "mongoose";

export default (connStr: string) => {
  const connect = async () => {
    try {
      await mongoose.connect(connStr);
      console.info(`successfully connected to ${connStr}`);
    } catch (error) {
      console.error(`error connected to ${connStr}`);
      process.exit(1);
    }
  };
  connect();

  mongoose.connection.on("disconnected", connect);
};
