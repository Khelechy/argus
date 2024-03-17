# Argus

Contributing to Argus Engine is pretty straightforward

- Fork the repo
- Make sure you have golang installed on your machine
- Make changes and raise a PR

## Contributing to building ARGUS Client Libraries ( Langauge Specific )

Building ARGUS client libraries invloves 3 straight simple steps.
- Connecting to ARGUS Engine
- Authenticating the connection
- Listening for events and messages

### Connecting to ARGUS Engine

Argus Engine leverages TCP/IP protocol is establish and maintain connections between the Engine and clients.

A simple TCP/IP dial/connect to the IP and Port is enough to do the trick.

```go
    import (
	    "net"
    )
    conn, err := net.Dial("tcp", "localhost:1337")

    defer conn.Close()
    ...    
```

```c#
    using System.Net;

    TcpClient client = new TcpClient();
    client.Connect("127.0.0.1", 1337);

    ...

    client.Close()
```

### Authenticating the connection

If authentication is turned on from the ARGUS Engine, it means clients connecting to the engine would have to authenticate the connection in order to receive the notification events.

The client has to send a prompt message to the engine immediately after connecting, if authentication is successful, the Engine adds it to the list of trusted clients which it would send events to.

The authentication message is a connection string in the formart `"<ArgusAuth>Username:Password</ArgusAuth>"`, where the `Username` and `Password` are the placeholders for the ARGUS Engine auth credentials.

```go
    connectionString := fmt.Sprintf("<ArgusAuth>%s:%s</ArgusAuth>", "testusername", "testpassword")

    sendAuthData(conn, connectionString)

    func sendAuthData(conn net.Conn, connectionString string) {
        data := []byte(connectionString)
        _, _ = conn.Write(data)
    }
```

```c#
    NetworkStream stream = client.GetStream();

    string authMsg = $"<ArgusAuth>{Username}:{Password}</ArgusAuth>";

    byte[] buffer = Encoding.ASCII.GetBytes(authMsg);
    stream.Write(buffer, 0, buffer.Length);

    stream.Close();
    client.Close();
```

### Listening for events and messages

In order to listen for events and messages you have to continously listen on the TCP stream for incoming data and then deserialize into an identifiable object ( for events), or log out messages for ordinary application messages.

```go
	buffer := make([]byte, 1024)
	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			argus.Errors <- err
		}

		data := string(buffer[:n])

		if len(data) > 0 {

			isJson, event, str := utils.IsJsonString(data)
	
			if isJson {
				// Push event to event channel
				argus.Events <- event
			} else {
	
				argus.Messages <- fmt.Sprintf("Received: %s\n", str)
			}
		}

	}
```

```c#
 var buffer = new byte[1024];
 int bytesRead;
 while (true)
 {
     bytesRead = _stream.Read(buffer, 0, buffer.Length);
     string response = Encoding.ASCII.GetString(buffer, 0, bytesRead);

     if (!string.IsNullOrEmpty(response))
     {

         var (isJson, argusEvent, message) = Helpers.IsJsonString(response);
         if (isJson)
         {
             OnRaiseCustomEvent(new ArgusEventArgs(argusEvent));
         }
         else
         {
             Console.WriteLine("Received: " + response);
         }
     }

 }
```

Return the fetched event to the user.

Note: The Json string expected from the ARGUS Engine is as below:

```json
	{
		"Action" : "string",
		"ActionDescription: "string",
		"Name": "string",
		"Timestamp": datetime
	}
```

