# QUIC-S

#file-synchronization #quic #server #go #golang #http3 #cobra

**QUIC-S** is a continuous file synchronization system based on QUIC protocol. 
This system is designed to synchronize files between multiple devices, to manage file versions, and to share files with other users.

**QUIC-S** consists of a server and a client, and the quic-go based quics-protocol is used for communication between the two.

This repository is a server repository.<br>
To see client: [quics-client](https://github.com/quic-s/quics-client)<br>
To see protocol: [quics-protocol](https://github.com/quic-s/quics-protocol)

> **NOTICE**: If you want to use this tool, you should use the client of QUIC-S. You can find the client in [here](https://github.com/quic-s/quics-client.git) 

[Features](#features) | [Getting Started](#getting-started) | [How to use](#how-to-use)| [Documentation](#documentation) | [Contribute](#contribute)

## Features

### 1. Manage client
Server can receive the registration request from the client. server saves the registered/connected client information to database badger. The saved client information is client uuid(created by client), IP address, etc.

### 2. Manage root directory per client
Server can manage root directory per client for synchronizing. The synchronization is performed in the registered root directory. When client request the root directory with registration, then server save the root directory to database badger, too. In addition the root directory can be registered one more.

### 3. Save & manage files synchronized from client
Similar to root directory, the file can be registered/saved to server. Requested file from client is updated with `latestHash` and `latestSyncTimestamp`. Server save latest files from client in their own directory (e.g., .quics/sync/${root-directory-name}/latest/*)

### 4. Manage & resolve conflict of file
If `LastUpdatedTimestamp` from client is larger than `LatestSyncTimestamp` from server, then any conflict could not be occurred. However, in the case of not above, conflict occurred.
When conflict occurs, then server makes a directory for managing conflict (e.g., .quics/sync/${root-directory-name}/conflict/*). The created conflict file can be removed after resolving conflict.
Server sends client with two options (client side, server side). Client chooses one option with two options, then sends chosen file with message to server. Server removes the conflict file, and create new file version/history about resolved file.

### 4. Save the history of file
Server manages all histories of all files. The history file is saved to directory (e.g., .quics/sync/${root-directory-name}/history/*). If the user wants, a file can be replaced with a previous file history.

> For more detail logic and implementation, please check [QUIC-S Docs](./docs/README.md)

## Getting Started

### 1. Docker

```Bash
docker run -it -d -e PASSWORD=passwordwhatyouwant -v /path/to/your/dir:/data --name quics -p 6120:6120 -p 6121:6121/udp -p 6122:6122/udp quics/quics
```

### 2. Local install
- 1. Download the latest version from [release page](https://github.com/quic-s/quics/releases)
- 2. Unpack the archive
- 3. Run `mv ./qis /usr/local/bin/qis`

### 3. Build from source
- 1. Install Go 1.21 or later.
- 2. Clone this repository.
     ```Bash
     git clone https://github.com/quic-s/quics.git
     ```
- 3. Run the command below in the root of the repository.
     ```Bash
     go mod download
     go build -o qis ./cmd
     ```

## How to use

### Environment variables

Environment variables are used to set initial server configuration. If you use docker, you can set environment variables with `-e` option.

| Name | Description | Default |
| - | - | - |
| QUICS_SERVER_ADDR | Server address | localhost |
| QUICS_SERVER_PORT | Legacy http port for Rest API server | 6120 |
| QUICS_SERVER_H3_PORT | Http/3 port for Rest API server | 6121 |
| QUICS_PASSWORD | Server password | password |
| QUICS_PORT | quics-protocol port for communication between server and client | 6122 |
| QUICS_CERT_NAME | Server certificate name for TLS | cert-quics.pem |
| QUICS_KEY_NAME | Server key name for TLS | key-quics.pem |

### CLI & REST API

**QUIC-S** currently supports CLI, but it all operates based on REST API. So you can use either cli or REST API for anything except starting the initial server process.

Below table is the list of commands and rest api path. 

> If you use docker, you meed to use `docker exec -it quics qis` or set alias `alias qis="docker exec -it quics qis"`.

| Tag | Command | Options | Description | Rest API |
| - | - | - | - | - |
| controller | `qis` | | root command meaning quic-s |
| controller | `qis` | `-h`, `--help` | show help |
| controller | `qis start` | | start rest server with default IP and port |
| controller | `qis start` | `--addr` string | start rest server with user-defined address |
| controller | `qis start` | `--port` string | start rest server with user-defined port for legacy http |
| controller | `qis start` | `--port3` string | start rest server with user-defined port for http/3 |
| controller | `qis run` | | run is a command that combines `qis start` and `qis listen` |
| controller | `qis run` | `--addr` string | start server with user-defined address |
| controller | `qis run` | `--port` string | start server with user-defined port for legacy http |
| controller | `qis run` | `--port3` string | start server with user-defined port for http/3 |
| controller | `qis listen` | | listen protocol | /api/v1/server/listen |
| controller | `qis stop` | | stop server | /api/v1/server/stop |
| config | `qis password set` | `--pw` string | change server password | /api/v1/server/password/set |
| config | `qis password reset` | | Reset server password | /api/v1/server/password/reset |
| log | `qis show` | | show various information |
| log | `qis show client` | `-i`, `--id` | show client information by key | /api/v1/server/logs/clients |
| log | `qis show client` | `-a`, `--all` | show all client information | /api/v1/server/logs/clients |
| log | `qis show dir` | `-i`, `--id` | show root directory information by key | /api/v1/server/logs/directories |
| log | `qis show dir` | `-a`, `--all` | show all root directory information | /api/v1/server/logs/directories |
| log | `qis show file` | `-i`, `--id` | show file information by key | /api/v1/server/logs/files |
| log | `qis show file` | `-a`, `--all` | show all files information | /api/v1/server/logs/files |
| log | `qis show history` | `-i`, `--id` | show history information by key  | /api/v1/server/logs/histories |
| log | `qis show history` | `-a`, `--all` | show all histories information | /api/v1/server/logs/histories |

## Documentation

For more detail logic and implementation, please check [QUIC-S Docs](./docs/README.md)

Also you can check [quics-client](https://github.com/quic-s/quics-client) for client side and [quics-protocol](https://github.com/quic-s/quics-protocol) for protocol.

## Contribute

**QUIC-S** is an open source project, and contributions of any kind are welcome and appreciated.

We also have a awesome plan to make **QUIC-S** better. Check [ROADMAP.md](https://github.com/quic-s/quics/blob/main/ROADMAP.md) will be helpful to understand our project's direction.

- To contribute, please read [CONTRIBUTING.md](https://github.com/quic-s/quics/blob/main/CONTRIBUTING.md)

- To report bugs or request features, please use the issue tracker. Before you do so, make sure you are running the latest version, and please do a quick search to see if the issue has already been reported.

- For more discussion, please join the [quics discord](https://discord.gg/HRtY7pNZz2)
