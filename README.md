# chat-encrypt

### How To Run

1. Make sure you have Makefile
2. For server, run `make serverx`
3. For client, run `make clientx`

### Directional Protocol

#### This protocol is sent by the server without user request

Note:

1. (output) = message from server

|         **Type**          | **Response**     |
| :-----------------------: | ---------------- |
| Other User Leave the Chat | LEAVE (output)\n |

<br>

### Bidirectional Protocol

#### User request to server and get the response back

Note:

1. (input) = input that you typed
2. (output) = message from server

| **Protocol** | **Request**        | **Success Response** | **Error Response** |
| ------------ | ------------------ | -------------------- | ------------------ |
| Registration | USERNAME (input)\n | RESPUSR (output)\n   | Error (output)\n   |
| Send Message | MESSAGE (input)\n  | RESPMSG (output)\n   | Error (output)\n   |
| Exit         | EXIT\n             | EXIT\n               | Error (output)\n   |

<br>

### Unknown Protocol

#### Response to unknown protocol

Note:

1. (output) = message from server

|     **Type**     | **Response**                                    |
| :--------------: | ----------------------------------------------- |
| Unknown Protocol | Unknown command. Please write correct command\n |
