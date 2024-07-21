# IO
- Utilities to load file, save file, zip file
- File Stream Writer
- File Stream Reader
- Implement ETL process for Data Processing, Business Intelligence

## ETL (Extract-Transform-Load)
Extract-Transform-Load (ETL) is a data integration process involving the extraction of data from various sources, transformation into a suitable format, and loading into a target database or data warehouse.
- Extracting data from various sources.
- Transforming the data into a suitable format/structure.
- Loading the transformed data into a target database or data warehouse.

## Batch processing
- [core-go/io](https://github.com/core-go/io) is designed for batch processing, enabling the development of complex batch applications. It supports operations such as reading, processing, and writing large volumes of data.
- [core-go/io](https://github.com/core-go/io) is not an ETL tool. It provides the necessary libraries for implementing ETL processes. It allows developers to create jobs that extract data from sources, transform it, and load it into destinations, effectively supporting ETL operations.

### Use Cases of [core-go/io](https://github.com/core-go/io) in ETL:
- <b>Data Migration</b>: Moving and transforming data from legacy systems to new systems.
- <b>Data Processing</b>: Handling large-scale data processing tasks like data cleansing and transformation
- <b>Data Warehousing</b>: Loading and transforming data into data warehouses.
- <b>Business Intelligence</b>: Transforming raw data into meaningful insights for decision-making, to provide valuable business insights and trends.

### Specific Use Cases of [core-go/io](https://github.com/core-go/io)
#### Export from database to file

  ![Export from database to file](https://cdn-images-1.medium.com/max/800/1*IEMXhQXJ0hWZBPL8q2jMNw.png)

##### Samples:
- [go-sql-export](https://github.com/project-samples/go-sql-export): export data from sql to fix-length or csv file.
- [go-hive-export](https://github.com/project-samples/go-hive-export): export data from hive to fix-length or csv file.
- [go-cassandra-export](https://github.com/project-samples/go-cassandra-export): export data from cassandra to fix-length or csv file.
- [go-mongo-export](https://github.com/project-samples/go-mongo-export): export data from mongo to fix-length or csv file.
- [go-firestore-export](https://github.com/project-samples/go-firestore-export): export data from firestore to fix-length or csv file.

#### Import from file to database

  ![Import from file to database](https://cdn-images-1.medium.com/max/800/1*rYaIdKGSd0HwZqZW7pMEiQ.png)
 
  - Detailed flow to import from file to database

    ![Import flow with data validation](https://cdn-images-1.medium.com/max/800/1*Y4QUN6QnfmJgaKigcNHbQA.png)

##### Samples:
- [go-sql-import](https://github.com/project-samples/go-sql-import): import data from fix-length or csv file to sql.
- [go-hive-import](https://github.com/project-samples/go-hive-import): import data from fix-length or csv file to hive.
- [go-cassandra-export](https://github.com/project-samples/go-cassandra-import): import data from fix-length or csv file to cassandra.
- [go-elasticsearch-import](https://github.com/project-samples/go-elasticsearch-import): import data from fix-length or csv file to elasticsearch.
- [go-mongo-export](https://github.com/project-samples/go-mongo-import): import data from fix-length or csv file to mongo.
- [go-firestore-export](https://github.com/project-samples/go-firestore-import): import data from fix-length or csv file to firestore.

##### Layer Architecture
- Popular for web development

![Layer Architecture](https://cdn-images-1.medium.com/max/800/1*JDYTlK00yg0IlUjZ9-sp7Q.png)

##### Hexagonal Architecture
- Suitable for Import Flow

![Hexagonal Architecture](https://cdn-images-1.medium.com/max/800/1*nMu5_jZJ1omzIB5VK5Lh-w.png)

##### Based on the flow, there are 4 main components (4 main ports):
- Reader, Validator, Transformer, Writer
##### Reader
Reader Adapter Sample: File Reader. We provide 2 file reader adapters:
- Delimiter (CSV format) File Reader
- Fix Length File Reader
##### Validator
- Validator Adapter Sample: Schema Validator
- We provide the Schema validator based on GOLANG Tags
##### Transformer
We provide 2 transformer adapters
- Delimiter Transformer (CSV)
- Fix Length Transformer
##### Writer
We provide many writer adapters:
- [SQL Writer](https://github.com/core-go/sql/blob/main/writer/writer.go): to insert or update data
- [SQL Inserter](https://github.com/core-go/sql/blob/main/writer/inserter.go): to insert data
- [SQL Updater](https://github.com/core-go/sql/blob/main/writer/updater.go): to update data


- [SQL Stream Writer](https://github.com/core-go/sql/blob/main/writer/stream_writer.go): to insert or update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- [SQL Stream Inserter](https://github.com/core-go/sql/blob/main/writer/stream_inserter.go): to insert data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush. Especially, we build 1 single SQL statement to improve the performance.
- [SQL Stream Updater](https://github.com/core-go/sql/blob/main/writer/stream_updater.go): to update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.


- [Mongo Writer](https://github.com/core-go/mongo/blob/main/writer/writer.go): to insert or update data
- [Mongo Inserter](https://github.com/core-go/mongo/blob/main/writer/inserter.go): to insert data
- [Mongo Updater](https://github.com/core-go/mongo/blob/main/writer/updater.go): to update data


- [Mongo Stream Writer](https://github.com/core-go/mongo/blob/main/batch/stream_writer.go): to insert or update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- [Mongo Stream Inserter](https://github.com/core-go/mongo/blob/main/batch/stream_inserter.go): to insert data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- [Mongo Stream Updater](https://github.com/core-go/mongo/blob/main/batch/stream_updater.go): to update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.


- [Elastic Search Writer](https://github.com/core-go/elasticsearch/blob/main/writer/writer.go): to insert or update data
- [Elastic Search Creator](https://github.com/core-go/elasticsearch/blob/main/writer/creator.go): to create data
- [Elastic Search Updater](https://github.com/core-go/elasticsearch/blob/main/writer/updater.go): to update data


- [Elastic Search Stream Writer](https://github.com/core-go/elasticsearch/blob/main/batch/stream_writer.go): to insert or update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- [Elastic Search Stream Creator](https://github.com/core-go/elasticsearch/blob/main/batch/stream_creator.go): to create data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.
- [Elastic Search Stream Updater](https://github.com/core-go/elasticsearch/blob/main/batch/stream_updater.go): to update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.


- [Firestore Writer](https://github.com/core-go/firestore/blob/main/writer/writer.go): to insert or update data
- [Firestore Updater](https://github.com/core-go/firestore/blob/main/writer/updater.go): to update data


- [Cassandra Writer](https://github.com/core-go/cassandra/blob/main/writer/writer.go): to insert or update data
- [Cassandra Inserter](https://github.com/core-go/cassandra/blob/main/writer/inserter.go): to insert data
- [Cassandra Updater](https://github.com/core-go/cassandra/blob/main/writer/updater.go): to update data


- [Hive Writer](https://github.com/core-go/hive/blob/main/writer/writer.go): to insert or update data
- [Hive Inserter](https://github.com/core-go/hive/blob/main/writer/inserter.go): to insert data
- [Hive Updater](httpshttps://github.com/core-go/hive/blob/main/writer/updater.go): to update data

- [Hive Stream Updater](https://github.com/core-go/hive/blob/main/batch/stream_writer.go): to update data. When you write data, it keeps the data in the buffer, it does not write data. It just writes data when flush.

## Summary
### File Reader
- File Stream Reader
- Delimiter (CSV format) File Reader
- Fix Length File Reader
### File Writer
- File Stream Writer
#### Delimiter (CSV format) Transformer
- Transform an object to Delimiter (CSV) format
- Transform an object to Fix Length format

## Appendix
### Import and export data for nodejs
#### Export data for nodejs:
##### Key features
- [onecore](https://www.npmjs.com/package/onecore): Standard interfaces for typescript to export data.
- [io-one](https://www.npmjs.com/package/io-one): File Stream Writer, to export data to CSV or fix-length files by stream.
##### Libraries to receive rows as stream, to export each record one by one:
- Postgres: [pg-exporter](https://www.npmjs.com/package/pg-exporter) to wrap [pg](https://www.npmjs.com/package/pg), [pg-query-stream](https://www.npmjs.com/package/pg-query-stream), [pg-promise](https://www.npmjs.com/package/pg-promise).
- Oracle: [oracle-core](https://www.npmjs.com/package/oracle-core) to wrap [oracledb](https://www.npmjs.com/package/oracledb).
- My SQL: [mysql2-core](https://www.npmjs.com/package/mysql2-core) to wrap [mysql2](https://www.npmjs.com/package/mysql2).
- MS SQL: [mssql-core](https://www.npmjs.com/package/mssql-core) to wrap [mssql](https://www.npmjs.com/package/mssql).
- SQLite: [sqlite3-core](https://www.npmjs.com/package/sqlite3-core) to wrap [sqlite3](https://www.npmjs.com/package/sqlite3).

##### Samples
- [oracle-export-sample](https://github.com/typescript-sample/oracle-export-sample): export data from Oracle to fix-length or csv file.
- [postgres-export-sample](https://github.com/typescript-sample/postgres-export-sample): export data from Posgres to fix-length or csv file.
- [mysql-export-sample](https://github.com/typescript-sample/mysql-export-sample): export data from My SQL to f11ix-length or csv file.
- [mssql-export-sample](https://github.com/typescript-sample/mssql-export-sample): export data from MS SQL to fix-length or csv file.

##### Import data for nodejs
###### Key features
- [onecore](https://www.npmjs.com/package/onecore): Standard interfaces for typescript to export data.
- [io-one](https://www.npmjs.com/package/io-one): File Stream Reader, to read CSV or fix-length files from files by stream.
- [xvalidators](https://www.npmjs.com/package/xvalidators): Validate data
- [import-service](https://www.npmjs.com/package/import-service): Implement import flow 
###### Libraries to write data to database
- [query-core](https://www.npmjs.com/package/query-core): Simple writer to insert, update, delete, insert batch for Postgres, MySQL, MS SQL
- Oracle: [oracle-core](https://www.npmjs.com/package/oracle-core) to wrap [oracledb](https://www.npmjs.com/package/oracledb), to build insert or update SQL statement, insert batch for Oracle.
- My SQL: [mysql2-core](https://www.npmjs.com/package/mysql2-core) to wrap [mysql2](https://www.npmjs.com/package/mysql2), to build insert or update SQL statement.
- MS SQL: [mssql-core](https://www.npmjs.com/package/mssql-core) to wrap [mssql](https://www.npmjs.com/package/mssql), to build insert or update SQL statement.
- SQLite: [sqlite3-core](https://www.npmjs.com/package/sqlite3-core) to wrap [sqlite3](https://www.npmjs.com/package/sqlite3), to build insert or update SQL statement.
- Mongo: [mongodb-extension](https://www.npmjs.com/package/mongodb-extension) to wrap [mongodb](https://www.npmjs.com/package/mongodb), to insert, update, upsert, insert batch, update batch, upsert batch.

##### Sample
- [import-sample](https://github.com/typescript-sample/import-sample): nodejs sample to import data from fix-length or csv file to sql (Oracle, Postgres, My SQL, MS SQL, SQLite)

## Installation
Please make sure to initialize a Go module before installing core-go/io:

```shell
go get -u github.com/core-go/io
```

Import:
```go
import "github.com/core-go/io"
```
