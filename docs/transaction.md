# Transaction

These are the transactions that can be performed over the network using quics-protocol. 

 
 > If you want to show image source, click [HERE](https://www.figma.com/file/8AIbhnR9RTnH8mVj9Qt8Qj/QUICS_Scenario?type=whiteboard&node-id=0%3A1&t=GEPrAAYiu4DPhx6g-1)

## Transaction List

* [Client Register](#client-register) 
* [Register Local Root Directory](#register-local-root-directory) 
* [Register Remote Root Directory](#register-remote-root-directory)
* [Please Sync](#please-sync)
* [Must Sync](#must-sync)
* [Force Sync](#force-sync)
* [Conflict](#conflict)
* [Full Scan](#full-scan)
* [Need Contents](#need-contents)
* [History Utils](#history-utils)
* [Sharing](#sharing)


## Client Register
![Client Register](https://github.com/quic-s/quics-client/assets/80394866/862b0eca-3fcd-47c9-a55e-25a90767c2d0)
Users must first register themselves as clients on QUIC-S. Register Client is a procedure similar to signing up.
If the registration is successful, you can use the synchronization service of QUIC-S in earnest.

STEPs are below :

1. The user sends a registration request to the server with the server's password
2. (1) If the server's password does not match, it fails
2. (2) If the server's password matches, check if the client is already registered
3. If the client is not already registered, it is registered as a new client on the server
 



## Register Local Root Directory
![Register Local Root Directory](https://github.com/quic-s/quics-client/assets/80394866/f133e2ca-7150-4ae3-becb-0a6deb89d858)

The client sends a request to register the root directory of the local directory to be synchronized to the server.
If the request is successfully completed, all files in the directory are synchronized.

STEPs are below :

1. The user selects the directory to be synchronized among the directories that have not yet been registered and sends a registration request
2. The server registers the root directory requested by the client. At this time, the password of the root directory sent by the client is also stored. This password is used when another client accesses the root directory
3. When a new root directory is registered, the server scans and synchronizes all directories it manages


## Register Remote Root Directory
![Registet Remote Root Directory](https://github.com/quic-s/quics-client/assets/80394866/8f3aa8ee-4452-4bb4-94d1-8b60c7395efa)

If you want to synchronize the root directory that another client has already registered, first check the list of root directories registered on the server. You can then select one and register the directory as your own root directory.


STEPs are below :

[**Remote Root Directory List**]
1. The user requests a list of root directories registered on the server
2. The server responds with a list of root directories registered on the server

[**Sync Remote Root Directory**]
1. The user selects one of the lists of root directories registered on the server and sends a registration request
2. The client compares the password of the root directory requested by the server with the password of the root directory registered on the server
3. If the password matches, the server adds the client to the synchronization target of the root directory
4. After processing the request, the server scans and synchronizes the entire directory managed by the server
5. The client registers the synchronization target directory sent by the server as its own root directory. From now on, you can receive changes to the directory in real time or send changes to the server

## Please Sync
![Please Sync](https://github.com/quic-s/quics-client/assets/80394866/0b442221-1ad8-4885-90c9-ffdc56e5a4ee)
The client sends its changes to the server. At this time, the changes mean changes such as file creation, modification, deletion, etc. The server receives and stores the changes.

STEPs are below :

1. The client sends the metadata of the changes it has made to the server
2. The server checks whether a conflict has occurred for the file
3. If a conflict occurs, the client is notified of the conflict and synchronization fails, and the process proceeds to the Conflict process
4. If no conflict occurs, the server requests the file contents
5. The client sends the file contents to the server
6. The server that received the file from the client performs Must Sync for the file

## Must Sync
![Must Sync](https://github.com/quic-s/quics-client/assets/80394866/3cd728b4-9dbc-4ac6-a84a-6a86cfbee91b)

Must Sync is the process of reflecting changes passed by other clients to the client. This is a request sent from the server to the client.

STEPs are below :

1. The server sends a synchronization request to the client that needs to synchronize the changed file
2. The client that received the request checks whether it is in a situation where it can synchronize the changed file locally
3. The client sends the server whether it can synchronize the request. If synchronization is not possible, the client requests Please Sync to the server and stops Must Sync
4. If the client is able to synchronize, the server sends the file to the client.

## Conflict
![Conflict](https://github.com/quic-s/quics-client/assets/80394866/0137fe11-d5e0-45e8-a072-4f612d3fc1bc)
Conflict can occur when multiple clients simultaneously modify the same file or when a client's internet connection is disconnected and synchronization is not performed. There are functions such as viewing the conflict list, downloading the contents, and resolving the conflict.

STEPs are below :

[**Show Conflict List**]
1. The client requests a list of conflicts from the server
2. The server responds with a list of conflicts
3. The client shows the conflict list received from the server to the user.

[**Download Conflict File**]
1. The client requests the number of files in which conflicts occurred from the server
2. The server returns the number of conflicts to the client. If the number of conflict files is 0, the process ends
3. If the number of conflicts is 0, the client ends the process.
4. If the number of conflicts is not 0, the server sends all the files in which conflicts occurred to the client
5. The client can download and check the contents of the conflict file.

[**Resolve Conflict by client**]
1. The client selects the file and candidate to be reflected in the conflict list and sends it to the server
2. The server processes the conflict based on the file selected by the client. The server also sends the final file selected by the client to the client. At this time, Force Sync is used.


## Force Sync
![Force Sync](https://github.com/quic-s/quics-client/assets/80394866/2b03eadc-88e7-460d-8530-e3a8120d776f)

Force Sync is a function that forces the client to synchronize and reflect the contents of the file stored on the server. This is a function that forces the client's contents to be overwritten when the contents of the file stored on the server are forcibly synchronized and reflected on the client.

STEPs are below :

1. The server sends a forced synchronization request to the clients to be synchronized
2. The client receives the forced synchronization request from the server. At this time, the contents of the file stored on the server are forcibly synchronized and reflected on the client.

## Full Scan
![Full Scan](https://github.com/quic-s/quics-client/assets/80394866/f1962b19-9304-4c97-bb0b-8279f8f61c34)

The server scans all files that are being synchronized and finds files that need to be synchronized after performing a full scan. It is mainly used to synchronize files that have been tracked for changes but have been missed. It can be done periodically or requested by the client.

STEPs are below :

1. The server collects metadata for the target client and all files for Full Scan and sends it to the client.
2. The client compares the information of the files in its local with the information received from the server. At this time, it also compares with the file metadata received from the local os
3. If they are not the same, the client goes to the process of requesting Please Sync to synchronize the file in the local.
4. All metadata including the information of the file that has gone to Please Sync is sent to the server, and Must Sync for the file that needs to be synchronized is started for the client.


## Need Contents
![Need Contents](https://github.com/quic-s/quics-client/assets/80394866/f337a5c6-ef7e-4998-8bf4-6b575dc9007a)

If the server needs the contents of a file that the client has, the server can get information from the client through the Need Contents process.


STEPs are below :

1. The server requests the client to send the contents of the file
2. The client sends the contents of the file to the server

## History Utils
![History](https://github.com/quic-s/quics-client/assets/80394866/409250ec-8080-4df9-97b5-be34a267dd52)

Server-managed files can be used for various functions using the file history stored and managed by the server. The functions are Rollback History, Show History, and Download History File.

STEPs are below :

[**Rollback History**]
1. First, the client sends a request to the server to roll back a specific file to a specific version
2. The server rolls back the file to the version requested by the client and synchronizes the rolled back file with all clients. 

[**Show History**]
1. The client sends the server a specific file and the number of histories to be viewed from **HEAD**
2. The server sends the client the history of a specific file.
At this time, the history includes information such as the version, modifier, and modification time of the file

[**Download History File**]
1. The client requests a specific version of a specific file from the server
2. The server sends the client the file of that version
3. The client can download and check the contents of the file received from the server

## Sharing
![Sharing](https://github.com/quic-s/quics-client/assets/80394866/24c6b187-b1e8-4274-bb31-2bc032554e6e)

QUIC-S provides a function to share files with third parties via links without using the application

STEPs are below :

[**Share File**]

1. Clients enter the file they want to share and the number of times they want to share
2. The server creates a link that can be shared as many times as the client wants and returns it to the client
3. The shared link is automatically deleted when the number of shares becomes 0

[**Stop Sharing**]
1. If Clients want to stop sharing via the link, request the server
2. The server stops sharing the file

