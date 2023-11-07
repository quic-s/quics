# Full Scan

This document describes the full scan process in QUIC-S system.

## What is Full Scan?

Full scan is a process of scanning all files in the directory and checking if there are any missing changes. 

This allows synchronization of files that were missed during the real-time synchronization process if the client was terminated or for some reason.

## When is Full Scan performed?

Full scan is performed when the client process is started or when the server requests it. The server requests a full scan from the client at preset time intervals. (Default: 5 minutes)

If user wants to perform full scan, user can use `qic rescan` command in client.

## Full Scan Process

To see full scan process, please check [transaction#FULLSCAN](fullscan.md#full-scan)