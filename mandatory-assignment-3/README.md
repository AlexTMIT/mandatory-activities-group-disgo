# ChittyChat 
Welcome to DisGo's ChittyChat Program! This README-file is here to provide you help when running the program, both as a server and as a client.
Please follow the guide step by step to run it smoothly.

## Guide to run the server
To run the server, first you have to make sure you're in the right folder: 'mandatory-assignment-3'. When that has been settled, do as follows:
1. Run the server: `go run server/main.go`.

And that's it! Now the server should be running.

## Guide to run the client
To run the client, assuming you are already in the right folder, do as follows:
1. Run the client: `go run client/client.go`.
2. Enter your name.
3. To enter a message, type `chat` followed by a space, followed by your message, e.g. `chat Hello!`. You can chat as many times as you want.
4. Once you are done with sending messages, and want to disconnect, simply type `leave`.

## Additional notes on the program
You can add more than one client to the server, in which each client asynchronously can send a message.
To disconnect the server, simply close the program, or press `ctrl` + `C`.
