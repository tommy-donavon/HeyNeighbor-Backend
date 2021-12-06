import express from 'express';
import { Chat, Message } from '../models/index.js';

const router = express.Router();

router.get('/past-msg', (req: express.Request, res: express.Response) => {
  let propC = req.query.property_code;
  let roomC = req.query.room;
  if (!propC || !roomC) {
    return res
      .status(404)
      .send(JSON.stringify({ message: 'please provide property code' }));
  }
  Chat.findOne({ property_code: <string>propC, room_name: <string>roomC })
    .lean()
    .populate({ path: 'messages' })
    .exec((err, chat) => {
      if (err)
        return res.status(500).send(JSON.stringify({ error: err.message }));
      if (!chat) return res.status(404).send(JSON.stringify({ messages: [] }));
      let output = chat.messages.sort((a,b) => a.time.getTime() - b.time.getTime())
      console.log(output)
      return res.status(200).send(JSON.stringify({ messages: output }));
    });
});

//TODO remove id from chat collection
router.delete('/msg', (req: express.Request, res: express.Response) => {
  let msgId = req.query.id;
  Message.findByIdAndDelete({ _id: msgId }, null, (err, result) => {
    if (err) return res.status(500).send(err);
    Chat.updateOne(
      { messages: result?._id },
      { $pullAll: { messages: [result?._id] } },
    ).exec((err, result) => {
      if (err) {
        console.error(err)
        return res.status(500).send(err);
      }
      return res.status(200);
    });
  });
});

export default router;
