# cassandra-fluency

## Architecture
- https://cassandra.apache.org/doc/stable/cassandra/architecture/index.html

## Unknown
- SSTable
- LSM Tree
    - https://youtu.be/MbwmMCu9ltg?si=0qIOk5sm-V-KEflC

## Terms
- Column Family
- Logical storage layout

| Primary Key | Columns |
|-------------|---------|
| 1 | (name, 'John'), (age, 30) |
| 2 | (name, 'Jane'), (age, 25) |
| 3 | (name, 'Doe'), (age, 35) |


- Attributes: keys_cache, rows_cache, preload_row_cache
- Primary Key: unique identifier for a row, no 2 rows can have the same primary key
    - Single column: Partition Key
    - Composite: Partition Key (how data is distributed) + Clustering Key (how data is ordered within a partition)
    - Partition key can be composite
    - There can be multiple clustering columns
- Keyspace ~ Database
- Table ~ Column Family
- Primary Key ~ Row key
- Analogy
```js
// generic
const columnFamily = {
    rowKey1: {
        colKey: colValue
    }
}

// specific
const clicksStream = {
    campaign1: {
        clicks: 100,
        date: '2024-02-14'
    },
    campaign2: {
        clicks: 200,
        date: '2024-02-14'
    }
}
```

| Primary Key | Columns |
|-------------|---------|
| campaign1 | (clicks, 100), (date, '2024-02-14') |
| campaign2 | (clicks, 200), (date, '2024-02-14') |

- Replication strategy
    - SimpleStrategy: 1 DC
    - NetworkTopologyStrategy: x DCs

## Features
- High write throughput
    - why?
        - Write to commit log first, then to memtable, then to sstable.
        - Memtable is in-memory, so write throughput is high.
- TTL
- Decentralized: No single point of failure, Every node is equal (any node can serve any request - be coordinator)
    - For read: Coordinator does not necessarily hold data, it just directs the request to the correct node. ?!?
- Wide column (column oriented): Store multiple columns in a single row
- Availability + Partition tolerance > Consistency / Performance
- Tunable consistency
    - Read consistency
    - Write consistency
## Warnings
- Do not use super column family as it loads all columns into memory, can lead to high memory usage and OOM.

## Operations
```sql
-- Create a keyspace
CREATE KEYSPACE catalog
    WITH REPLICATION = {
        'class' : 'SimpleStrategy', -- Fault-tolerant
        'replication_factor' : 3 -- 1 master, 2 replicas
    };

-- use keyspace
USE catalog;

-- Create column family
CREATE COLUMNFAMILY product (
    productId varchar,
    title text,
    brand varchar,
    publisher varchar,
    length int,
    width int,
    height int,
    PRIMARY KEY (productId)
);
    
-- describe all column families
DESCRIBE COLUMNFAMILIES;

-- describe column family
DESCRIBE product;

-- alter column family
ALTER COLUMNFAMILY product ADD price double;

-- insert data
INSERT INTO product (productId, title, brand, publisher, length, width, height, price) 
VALUES ('1', 'Product 1', 'Brand 1', 'Publisher 1', 100, 100, 100, 100);
-- insert data with TTL 10 seconds
INSERT INTO product (productId, title, brand, publisher, length, width, height, price) 
VALUES ('2', 'Product 1', 'Brand 1', 'Publisher 1', 100, 100, 100, 100) USING TTL 10;

-- select data
SELECT * FROM product;

-- create productviewcount column family
CREATE COLUMNFAMILY productviewcount (
    productId varchar,
    viewCount counter,
    PRIMARY KEY (productId)
);

-- increment view count
UPDATE productviewcount SET viewCount = viewCount + 1 WHERE productId = '1';
```

- Advanced
    - Collection data types
        - Set
        - List
        - Map
    - Counter

- Ring topology (consistent hashing)
    - Each node is assigned a initial_token, which determines the range of tokens that the node is responsible for.
    - Partitioner: 
        - hash(partitionKey) -> token
            - All rows with the same partition key will be assigned to the same node.
            - Entire row is stored on the node that is responsible for the token.
        - Strategy:
            - RandomPartitioner: md5
            - Murmur3Partitioner: murmur3 hash (default)

- Restrictions
    - `select * from table where col = ?`: all columns in the partition key must be specified in the where clause.
- Secondary Index
    - Non-primary key columns can be indexed.
    - Low cardinality columns are good candidates for secondary indexing.
    - Use separate column families for secondary indexing, no replication.
    - Query is sent to all nodes
    - Range queries: must use ALLOW FILTERING

- To which nodes, replicas are sent?
    - Token
    - Replication placement strategy

- Write consistency
- Read consistency
    - Snitch
    - level 1, level all, level quorum, level local quorum
    - quorum: resp from node A + hash(resp) from node B
    - fix inconistency
        - read repair
     
## Ucase: Adtech
- https://engineeringblog.yelp.com/2016/08/how-we-scaled-our-ad-analytics-with-cassandra.html
- Timeseries: https://youtu.be/5vWvukzk9Z0?si=Xr3yZoPq-yj9sgI0
- Timeseries: https://www.youtube.com/watch?v=U7oaBDlXvhc
