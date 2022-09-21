# igbodb
## Object oriented database

## Smart RAFT based consensus
### Best machine is always the master. Does it improve performance?
### If multiple machines one is used as a backpressure proxy? What is the performance of that.
## Reactive first

## Eventually consistent

## Commit listeners
Can a client listen on a transaction commit in a reactive mode?

## Sharding 

## Indexing 

## Query language support

## Leveraging RocksDB capabilities
https://github.com/facebook/rocksdb
Using grooksdb as a wrapper for rocksdb
https://github.com/linxGnu/grocksdb


## Sharding
Efficent sharding mechanism based on Cocroach DB. Need more research on best way to do it.

## Report only nodes
Querying can be slow in some scenarios where multiple joins and traversals of object trees are required. Indexing is managed as a separated module and can be deployed as a separated node or group of nodes.


## Implementation details
Client do any translation work required
Client always adhere to simplicity for the language it is used in

Client to send the data structure of a object every time it changes or every time? Every time to start.

Data structures to be used?
UUID encoded for uniquenes of an object and encoding for shorting the size and save disk space

every record has to be an object

object as attribute of another object will be linked by the uuid.

arrays will be its own object.

Implement object extensions.

Basic represenation of an object

-- Object A
--- Attribute A - S (String)
--- Attribute B - I (Integer)
--- Attribute C - A (Array) of -> Object B Object ref to Array Z
--- Attribute D - O (Object) B -> Object ref to Object B

-- Object B
--- Attribute X - String
--- Attribute Y - Blob

-- Object Array/Map Z
--- key String
--- value Object ref


Map will be like array but instead of an index as key it will have a given key



