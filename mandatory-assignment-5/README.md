# A Distributed Auction System

Welcome to DisGo's version of a distributed auction system! This README-file is here to provide you assistance when running the program. There are two runnable files:
* `process/process.go`
* `client/client.go`

## To run the process(es)
To run `process.go`, simply direct your terminal to the `mandatory-assignment-5` directory and type: 
```
go run process/process.go
```
This leads you to a question of what port you want your server to run on. Currently, the program only accepts two types of inputs:
* Localhost on port 50051
* Localhost on port 50052

You can choose to run either by typing `51` or `52` in your terminal, or run both, by opening two terminals and run one port in each.
There! Now your server should be up and running.

## To run the client(s)
When your server(s) is finally running, you can create a client. To run `client.go` simply direct your terminal to the `mandatory-assignment-5` directory and type: 
```
go run client/client.go
```
Here, you can choose to run one or more client. Every client will automatically connect to the two ports.

Type in your name to enter the auction. When you have written your name, you will enter the auction.

To bid in the auction, type
```
bid <amount>
```
and to query the auction, type
```
query
```
You may continue to bid until the auction has ended, and query however many times you like!
