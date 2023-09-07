# ROADMAP

This document defines a high level roadmap for QUIC-S development and upcoming releases. The features and themes included in each milestone are optimistic in the sense that many do not have clear owners yet. Community and contributor involvement is vital for successfully implementing all desired items for each release. 

We hope that the items listed below will inspire further engagement from the community to keep QUIC-S progressing and shipping exciting and valuable features.

Any dates listed below and the specific issues that will ship in a given milestone are subject to change but should give a general idea of what we are planning. 

## V1.0

### Objective time: In 2023

V1.0 focuses on the core features of QUIC-S. The goal is to have a stable and performant implementation of the QUIC-S that can be used in production environments.

### The key features of V1.0

- Real time synchronization

The tool should be able to synchronize files in real time. This means that when a file is changed on one computer, the change should be propagated to the other computers as soon as possible.

- Large file support

The tool should be able to handle large files. This means that the tool should be able to synchronize files that are several gigabytes in size.

- Conflict resolution

The tool should be able to resolve conflicts between files. This means that if a file is changed on two computers at the same time, the tool should be able to resolve the conflict by asking the user which version of the file to keep.

The standard of time of files is logical time. Logical time is added from 1, when the file have change events. So, the file with the latest logical time is the latest file.

- History management

The tool should be able to manage the history of files. This means that the tool should be able to keep track of all the changes that have been made to a file, and should be able to revert to a previous version of the file if necessary.

## V1.5

### Objective time: 2024 First Half

This version aims to improve the user experience of the file synchronization program completed in 1.0. 

### The key features of V1.5

- Optimized file synchronization

Measure the speed of synchronization by measuring the performance of previous versions, etc. and see if there are any areas for optimization. This will be an ongoing consideration for each version.

- Add a GUI

Add a GUI to the file synchronization program to make it easier to use. Client CLI is already use REST API, so it is ready to add GUI.

- Add support of file transfer protocols (FTP, SFTP, WebDAV, etc.)

Add support of file transfer protocols (FTP, SFTP, WebDAV, etc.) to the file synchronization program to access files remotely. 

- Sharing files with other users

Add support of sharing files with other users that not use QUIC-S. This need to add a feature that can share files with other users by sending a link to the file. Also, it need to add access control to the file.

## QUIC-S Next

### Objective time: After V1.5

After V1.5, QUIC-S aims to go beyond basic file synchronization and sharing and add more professional and advanced features.

- Add Mobile app (Android, iOS)

Add Mobile client app (Android, iOS) to the file synchronization program to synchronize files on mobile devices. 

- Add distributed mode

Add distributed mode to the quics server to synchronize files between multiple servers. This increases availability by ensuring that services can continue to be provided even if one server fails, and increases reliability by storing files in a distributed system.

This is a highly technical task and will require research into things like the RAFT algorithm.

- Add support of deploying on Kubernetes

Add support of deploying on Kubernetes to the quics server to make it easier to deploy the quics server on Kubernetes. This will need to add helm chart for easy deployment.

