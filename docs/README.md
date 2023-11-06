# QUIC-S Docs

QUIC-S Docs are documents for how QUIC-S works. Write a description of the logic and implementation that is too detailed to be written in README.

## Table of Contents

- [System Architecture](system-architecture.md)
- [Transaction](transaction.md)
- [Conflict](conflict.md)
- [History](history.md)
- [FullScan](fullscan.md)

## How to write

For significant scope and complex new features, it is recommended to write a Document before starting any implementation work. On the other hand, we don't need to documentation for small, simple features and bug fixes.

Writing a document for big features has many advantages:

- It helps new visitors or contributors understand the inner workings or the architecture of the project.
- We can agree with the community before code is written that could waste effort in the wrong direction.

While working on your document, writing code to prototype your functionality may be useful to refine your approach.

Authoring document is also proceeded in the same [contribution flow](../CONTRIBUTING.md) as normal Pull Request such as function implementation or bug fixing.