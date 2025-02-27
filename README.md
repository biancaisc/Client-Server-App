# Client Server App 

Client - Server type system (multiple clients and one server), where the data processing method on the Server is done concurrently - using goroutines. 

## Description
There is a [configuration file](configurations.json) in which there are initial parameters of the program (maximum size of the data array the client can send, maximum number of goroutines)

If the number of active goroutines reaches the limit, the server waits before accepting new clients.

There are messages between client and server like: 
- Client `<Name>` Connected.
- Client `<Name>` made data request: `<date>`.
- Server has received the request.
- Server processes the data.
- Server sends `<response>` to client.
- Client `<Name>` received the response: `<response>`. 

## Example Usage
The system is tested with some basic exercises where the client sends an array of data to the server, the server processes the array and sends its response back to the client.
