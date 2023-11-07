# Conflict

This document describes the conflict state and how to handle it.

## What is Conflict?

Conflict refers to inconsistencies in the contents of a file during file synchronization, either because multiple clients are modifying the same place on a particular file at the same time, or because one client's internet connection is disconnected and synchronization cannot proceed.

## Conflict Detection

Conflict detection is a process of detecting conflicts in a file. The server detects conflicts by comparing the timestamp and hash of the file when the client sends the file to the server while PLEASESYNC transaction.

The conditional expression for detecting conflicts is as follows:

```
serverFile.LatestSyncTimestamp < clientFile.LastUpdateTimestamp && serverFile.LatestHash == clientFile.LastSyncHash
```

### Conflicted File

If the conditional expression is true, the file is conflicted. The server creates a conflict data and stores it in the database.

At this point, no more files will be updated and any conflicting files will be stored in a directory named `<syncRootDir>.conflict`. They are stored in the form `filename_<modified client uuid>` and wait for the user to resolve the conflict.

### Conflict Data

```go
type Conflict struct {
	AfterPath    string
	StagingFiles map[string]FileHistory
}
```

The conflict data is a struct that stores in the database. It contains the path of conflicted file and the data of out-of-sync files that are stored in the staging area(conflict directory).

## Conflict Resolution

Conflict resolution is a process of resolving conflicts in a file. The user can resolve the conflict by selecting one of the candidate files stored in the conflict directory.

After the user selects the file, the server sends the file as FORCESYNC transaction to the client. The client receives the file and replaces the conflicted file with the selected file.
