# A.R.G.U.S

ARGUS is a versatile and lightweight file/folder watcher designed for seamless integration into any environment. 
It provides real-time notification events for file and folder operations, allowing connected clients to stay updated on changes within the watched directories. 
Leveraging the power of Go's concurrency and the cross-platform support of fsnotify, ARGUS ensures robust performance across various operating systems.

## Available Client Libraries
- [Argus Go](https://github.com/Khelechy/argus-go)
- [Argus .NET](https://github.com/Khelechy/argus-dotnet)

## Dependencies

- golang +v1.19
- fsnotify

## Features

- Real-time Notification: Receive instant notifications for file and folder operations, including creations, deletions, modifications, change modes(chmod) and renames.

- TCP/IP Communication: Connect to ARGUS via TCP/IP to receive notifications, enabling seamless integration with a wide range of applications and services.

- Multiple Clients: ARGUS supports multiple connected listening clients with no limit.

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
go build
```

Run argus executable with the command
```sh
./argus
```

Note: If you are in the current working directory of argus cloned project you dont need to pass a config file as flag. But if you are not, you need to pass a config file path as a command line argument.

```sh
./argus -config=/path/to/config.json
```

Sample config file
```json
{
    "server": {
        "host": "localhost",
        "port": "1337",
        "username": "testuser",
        "password": "testpassword"
    },
    "watch": [
        {
            "path": "C:/Users/PFY-102.PFY-102/source/repos/Mine/argus/testfolder2",
            "isFolder": true,
            "watchRecursively": true
        },
        {
            "path": "C:/Users/PFY-102.PFY-102/source/repos/Mine/argus/testfolder",
            "isFolder": true,
            "watchRecursively": false
        },
        {
            "path": "C:/Users/PFY-102.PFY-102/source/repos/Mine/argus/file.txt",
            "isFolder": false,
            "watchRecursively": false
        }
    ]
}
```     


## Contribution
Feel like something is missing? Fork the repo and send a PR.

Encountered a bug? Fork the repo and send a PR.

Alternatively, open an issue and we'll get to it as soon as we can.

### Building Target Language Libraries

To contribute to building the ARGUS client libraries and conforming with the underlying standards and procedures please read the [contributing docs](https://github.com/Khelechy/argus/blob/main/CONTRIBUTING.md)
