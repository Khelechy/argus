# A.R.G.U.S

ARGUS is a versatile and lightweight file/folder watcher designed for seamless integration into any environment. 
It provides real-time notification events for file and folder operations, allowing connected clients to stay updated on changes within the watched directories. 
Leveraging the power of Go's concurrency and the cross-platform support of fsnotify, ARGUS ensures robust performance across various operating systems.

## Dependencies

- golang +v1.19
- fsnotify

## Features

- Real-time Notification: Receive instant notifications for file and folder operations, including creations, deletions, modifications, change modes(chmod) and renames.

- TCP/IP Communication: Connect to ARGUS via TCP/IP to receive notifications, enabling seamless integration with a wide range of applications and services.

- Authentication Support: Ensure secure communication by authenticating clients connecting to ARGUS, providing an additional layer of protection for sensitive data.

- Multiple File/Folder Watching: Watch multiple files and folders concurrently, allowing users to monitor various locations simultaneously.

- Recursive Folder Watching: Enable recursive watching to track changes within nested directories, providing comprehensive coverage of file system activity.

- Cross-platform Compatibility: Utilize ARGUS across different operating systems, thanks to its dependency on fsnotify, which offers consistent performance and behavior regardless of the platform.

- Flexible Configuration: Includes a default `config.json` file in the root of the project for configurations samples and authentication credentials, (Leave username and password empty if you dont need auth).


# Getting Started:

To clone and build ARGUS, follow these simple steps:

Clone repo to your local directory with the command
```sh
git clone https://github.com/khelechy/argus.git
```

Build an executable with the command
```sh
cd argus
go build -o argus main.go
```

Run argus executable with the command
```sh
./argus
```

Note: If you are in the current working directory of argus cloned project you dont need to pass a config file as flag. But if you are not, you need to pass a config file path as a command line argument.

```sh
./argus -config=/path/to/config.json
```
