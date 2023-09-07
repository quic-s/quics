# quics
quics is a server for the QUIC-S. It is continuous file synchronization tool based to the QUIC protocol.

#quic #server #go #golang #http3 #cobra

## Features
### 1. Manage client
Server can receive the registration request from the client. server saves the registered/connected client information to database badger. The saved client information is client uuid(created by client), IP address, etc.

### 2. Manage root directory per client
Server can manage root directory per client for synchronizing. The synchronization is performed in the registered root directory. When client request the root directory with registration, then server save the root directory to database badger, too. In addition the root directory can be registered one more.

### 3. Save & manage files synchronized by client
Similar to root directory, the file can be registered/saved to server. Requested file from client is updated with `latestHash` and `latestSyncTimestamp`. Server save latest files from client in their own directory (e.g., .quics/sync/${root-directory-name}/latest/*)

### 4. Manage & resolve conflict of file
If `LastUpdatedTimestamp` from client is larger than `LatestSyncTimestamp` from server, then any conflict could not be occurred. However, in the case of not above, conflict occurred.
When conflict occurs, then server makes a directory for managing conflict (e.g., .quics/sync/${root-directory-name}/conflict/*). The created conflict file can be removed after resolving conflict.
Server sends client with two options (client side, server side). Client chooses one option with two options, then sends chosen file with message to server. Server removes the conflict file, and create new file version/history about resolved file.

### 4. Save the history of file
Server manages all histories of all files. The history file is saved to directory (e.g., .quics/sync/${root-directory-name}/history/*). If the user wants, a file can be replaced with a previous file history.

## Getting Started
### 1. Docker
```Bash
docker run -it --rm -v /path/to/your/dir:/data chromato99/quics
```

### 2. Local install
- 1. Download the latest version from this repository.
- 2. Unpack the archive
- 3. Run `go build ./cmd`

### 3. Build from source
- 1. Install Go 1.21 or later.
- 2. Clone this repository.
     ```Bash
     git clone https://github.com/quic-s/quics
     ```
- 3. Run the command below in the root of the repository.
     ```Bash
     go build ./cmd
     ```

## Contribute
- To report bugs or request features, please use the issue tracker. Before you do so, make sure you are running the latest version, and please do a quick search to see if the issue has already been reported.
- For more discussion, please join the [quics discord](https://discord.gg/HRtY7pNZz2)
