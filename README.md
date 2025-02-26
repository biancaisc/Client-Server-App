# Client Server App 

Client - Server type system (multiple clients and one server), where the data processing method on the Server is done concurrently - using go routines. 

## Description
There is a [configuration file](configurations.json) in which there are initial parameters of the program (how many elements of the data array the client can send, how many 
times a go routine can be called).

There are messages between client and server like: 
- Client <Name> Connected.
- Client <Name> made data request: <date>.
- Server has received the request.
- Server processes the data.
- Server sends <response> to client.
- Client <Name> received the response: <response>. 
