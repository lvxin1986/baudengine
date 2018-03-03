# baud

Baud is a multi-model storage system with reliability, scalability, and flexibility. 

## Data Model

Field, Object, Class, Edge, Space

each object has a field as the 'partition key' and an internal 'Object ID' (OID) generated by the system. 

each field can be indexed. 

## Cluster Management

runs on JDOS/Kubernetes across servers, racks, cells and regions

## Partitioning

space --> partition

partition ID range

## Replication

raft

async

filtered

## Components

there are several kinds of nodes in a Baud cluster: router, toposerver, objectserver, blobserver, which all speak with HTTP+JSON

## Search

strong indexing and search capability

## Graph

src --edge--> dst

## BLOBs

## File Extents


