# -- TCP/IP Simulator in Go -- Question Assignment
We have implemented two solutions to assignment (1)[easy]. However, the answers below pertain to the solution called "tcp-ip-simulator", which does not use sequence numbers, but rather enums, and uses timeouts. The alternative/second solutions has sequence numbers, but not timeouts. Lastly, the tcp-ip-simulator.go can run with a multitude of servers and clients, whereas the alternative solution has one server and one client. 

### a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?
Answer: The following packages, and the reason for their use, have been implemented:
- fmt: for printing states
- math/rand: for randomising the sequence numbers that the client / server sends
- time: for measuring time in case of time outs

In order for the client and the server to be communicating, we use channels to transmit data.

### b) Does your implementation use threads or processes? Why is it not realistic to use threads?
Answer: Yes, to simulate the TCP/IP handshake, we've used threads - a client and a server - to communicate with each other through the channels. 
Although, it is not realistic, considering threads use the same process, whereas a client and a server in reality use different processes.

### c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
Answer: In the case of messages arriving out of order, we would re-order the messages by implementing sequence numbers.

### d) In case messages can be delayed or lost, how does your implementation handle message loss?
Answer: In case SYN, SYN_ACK, or ACK have been lost of delayed more than 5 seconds, the client or server resends the message. There is, however, no cap to how many times a message may be lost. 

### e) Why is the 3-way handshake important?
Answer: To establish connection, to synchronise their sequence numbers, to make sure both the client and server are ready and have a stable and reliable connection, and that they're both able to remain connected.

________

# Assignment Description (excluding the questions, because they're right above)
Implement the TCP/IP protocol in Go. Your implementation has to be a simulation of the protocol seen in class (see slides).

There are different levels that you can work on. In order to pass, you need to implement at least (1) or (2).

(1)[Easy] Implement the TCP/IP Handshake using threads. This is not realistic (since the protocol should run across a network) but your implementation needs to show that you have a good understanding of the protocol. 

(2)[Hard] Implement a TCP/IP Handshake using the net package.

(3)[Medium] Implement a forwarder process/thread that simulates the middleware, where messages can be delayed or lost. All messages must go through the forwarder.     
