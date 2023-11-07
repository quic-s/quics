# System Architecture

**QUIC-S** is a continuous file synchronization system based on QUIC protocol. 
This system is designed to synchronize files between multiple devices, to manage file versions, and to share files with other users.

This document is written for describe the system architecture and its components.

## Table of Contents

* [Overall Architecture](#overall-architecture)
* [quics-protocol](#quics-protocol)
    * [Why develop a new protocol?](#why-develop-a-new-protocol)
    * [Features of quics-protocol](#features-of-quics-protocol)
* [Restful API](#restful-api)
    * [API Structure](#api-structure)
    * [Implementation](#implementation)
* [quics and quics-client](#quics-and-quics-client)
    * [Code Structure](#code-structure)
    * [Network](#network)
    * [File System](#file-system)
    * [Database](#database)
    * [Command Line Interface](#command-line-interface)

## Overall Architecture

![architecture](https://github.com/quic-s/quics/assets/20539422/c41053e8-3786-4df9-b426-3e8f4041eebb)

QUIC-S is a file synchronization system that allows a centralized server to synchronize files in multiple client directories in real time. Therefore, it has a client-server structure for file synchronization and uses the [quics-protocol](https://github.com/quic-s/quics-protocol) to communicate between them.

In this architecture, the server is a [quics](https://github.com/quic-s/quics) and the client is a [quics-client](https://github.com/quic-s/quics-client).

However, the unique thing is that each client and server also has the role of an http server for the Restful API, so the client and server can be controlled through this Restful API.

Therefore, the CLI commands (`qic` and `qis`) to control the client and server also use http client internally to interact with the client and server respectively, except for the command that initially starts the process.

This allows us to build integrations with other systems without necessarily using the CLI, and we have a web-based GUI in mind for the future.

## quics-protocol

The quics-protocol was developed as part of the QUIC-S project for communication between clients and servers in QUIC-S.

It is based on the QUIC protocol (using [quic-go](https://github.com/quic-go/quic-go) that implements the QUIC protocol in Go), and the protocol is designed to sending and receiving byte messages and files.

> For more detail structure, please check [quics-protocol](https://github.com/quic-s/quics-protocol)

### Why develop a new protocol?

The reason for developing a new protocol is that the existing file synchronization protocols are not suitable for QUIC-S.

Existing file synchronization protocols are based on TCP, but classic TCP does not have built-in stream multiplexing, compression, and security like QUIC, and these features must be implemented at the application layer. 

Therefore, we thought it would be a good idea to implement a new lightweight and simple protocol that makes full use of the standardized QUIC features.

### Features of quics-protocol

#### 1. Transaction

The quics-protocol is designed to allow servers and clients to send and receive data in units of communication called transactions. 

Each transaction creates a single stream in QUIC and can send and receive data multiple times before closing.  Also, since independent streams are created, transactions can communicate in parallel without any head of line issues.

#### 2. Handler

The handler is an object created and used internally by the quics-protocol object, which connects the transaction to a predefined callback function based on the transactionName.

This is similar to a mux in HTTP and uses a method called RecvTransactionHandleFunc to predefine the callback function. 

#### 3. FileInfo

The quics-protocol sends fileInfo along with the file when it sends a file, so that it can be sent to a remote location identical to the original (permissions, modification time, name, etc.). 

This allows file synchronization systems to maintain the exact same file on each device.

## Restful API

quics and quics-client have a Restful API for controlling the program. 

This API is used internally by the CLI commands, but it can also be used to control the client and server without using the CLI.
It also allows us to build integrations with other systems without necessarily using the CLI.

In addition, we are planning to develop a web-based GUI in the future, and this API will be used to communicate with the client and server.

### API Structure

The detail of API structure is in each repository's README.md.

* [quics](https://github.com/quic-s/quics)
* [quics-client](https://github.com/quic-s/quics-client)

### Implementation

The Restful API is implemented using Go's standard library `net/http` and `gorilla/mux`.

The `net/http` package is a package that implements the HTTP protocol, and the `gorilla/mux` package is a package that implements the router for the HTTP protocol.

But also we support the http/3 protocol, so we use the `quic-go` package that implements the QUIC protocol in Go.

Therefore, each process uses two ports to simultaneously support the existing TCP-based http and http/3. The default port for http is 6120, and the default port for http/3 is 6121.

## quics and quics-client

The quics and quics-client are the core components of the QUIC-S system.
quics is a server that manages file synchronization, and quics-client is a client that synchronizes files with the server.

### Code Structure

Hexagonal architecture is a model of designing software applications around domain logic to isolate it from external factors. The domain logic is specified in a business core, which we’ll call the inside part, with the rest being outside parts.

The advantage of hexagonal architecture is that it allows us to isolate our core business logic from the external factors that may affect it, such as file IO, database, network, etc. By using ports and adapters, we can decouple our domain model from the outside world and make it easier to test, debug, maintain, and change

### Network

In quics and quics-client, the network package implements the adapters needed for networking. Current quics systems use the quics-protocol and http for RESTful APIs, so these adapters are implemented. 

#### quics-protocol

The quics-protocol implements the communication between client and server within the quics system. It allows transactions such as client registration, syncRootDir registration, PleaseSync, MustSync, and FullScan to work. 

The client and server must receive or make requests depending on the transaction, so there is a primary adapter for receiving requests and a secondary adapter for making requests.

For details of transactions, please refer to the [Transaction](transaction.md) document.

#### http

HTTP is used for communication to externally control the QUICS system. This includes the command line interface used by users. 

The interface is implemented as a RESTful API, and is available through the http and http/3 protocols. The CLIs for controlling the server and client, `qis` and `qic`, use http/3.

The client and server have a primary adapter implemented to accept these requests. 

For details of the API, please refer to each repository's README.md.

* [quics](https://github.com/quic-s/quics)
* [quics-client](https://github.com/quic-s/quics-client)

### File System

The file system package implements the adapters needed for file system operations. The quics system uses the file system to manage files, so these adapters are implemented.

#### Server

The server accesses the file system by default to store files received from clients for synchronization and to transfer them to other clients, so an adapter for this is implemented using Go's OS package.

These include read/write of files, read/write of history, and read/write of concurrent files. Accessing the same file at the same time can cause errors, so mutex locks are applied based on the file's path.

#### Client

The client needs to be able to listen for CREATE, WRITE, RENAME, and DELETE events from the file system so that synchronization of files can be triggered. 

For this purpose, a library called fsnotify is implemented to receive events from the operating system's file system.
And an adapter is implemented to perform read and write operations to synchronize files when there is an event or transaction.

### Database

The QUICS system uses the [badger](https://github.com/dgraph-io/badger) library as a database to store and manage data.

BadgerDB is a fast key-value store written in pure Go that uses the [LSM tree](https://en.wikipedia.org/wiki/Log-structured_merge-tree) data structure. 

It is suitable for the quics system because it can perform fast and parallel transactions, is optimized for SSDs, and is embeddable.
This is used to store information from clients and files, and to handle history and conflicts.

The struct of the data stored in the database is defined in the types package, and the adapter for accessing the database is implemented in the repository package.

### Command Line Interface

quics and quics-client have a command line interface for controlling the program named `qis` and `qic` respectively.

This CLI is implemented using the [cobra](https://github.com/spf13/cobra) library, and the command is implemented in the cmd package.

The CLI uses the http client to communicate with the server and client, and the http client is implemented internally.

For details of the CLI commands, please refer to each repository's README.md.

* [quics](https://github.com/quic-s/quics)
* [quics-client](https://github.com/quic-s/quics-client)
