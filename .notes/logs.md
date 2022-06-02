# Logs

## Terms

**Log**:

- the abstraction that ties all the segments together
- an append-only sequence of records
- records are appended to the end
- typically read top to bottom, oldest to newest
- any data can be logged
- is split into a list of segments for disk space management
- always orders the records by time
- always indexes each record by its offset and time created
- always has one special segment known as the active segment

**Record**:

- the data stored in our log
- reading a record given its offset is a two-step process:
    - first you get the entry from the index file for the record, which tells you the position of the record in the store file
    - second you read the record at that position in the store file

**Offset**:

- a unique sequential number assigned to a record that acts like an ID
- assigned by the log when appending a record to the log

**Segment**:

- the abstraction that ties a store and an index together
- are used to deal with limited disk space
- when the log grows too big, we free up disk space by deleting old segments whose data we've already processed or archived
- comprises of a store file and an index file

**Active Segment**:

- special segment among the list of segments in a log
- the only segment we actively write to
- when we've filled the segment, we create a new segment and make it the active segment

**Store File**:

- the file we store records in
- stores the record data
- records are continually appended to this file

**Index File**:

- the file we store index entries in
- indexes each record in the store file
- speeds up reads because it maps record offsets to their position in the store file
- is much smaller than the store file that stores all your record data
- small enough that we can memory-map them and make operations on the files as fast as operating on in-memory data
- requires two fields: 
    - the offset
    - the stored position of the record


## How Logs Work

A log is an append-only sequence of records. You append records to the end of the log,
and you typically read top to bottom, oldest to newest - similar to running `tail -f`
on a file. You can log any data. 

People have historically used the term logs to refer to lines of text meant for humans to read, 
but that's changed as more people use log systems where their "logs" are binary encoded messages
meant for other programs to read.

When you append a record to a log, the log assigns the record a unique and sequential offset
number that acts like the ID for that record. A log is like a table that always orders the
records by time and indexes each record by its offset and time created.

Concrete implementations of logs have to deal with us not having disks with infinite space,
which means we can't append to the same file forever. So we split the log into a list of
segments. When the log grows too big, we free up disk space by deleting old segments whose
data we've already processed or archived. This cleaning up of old segments can run in a
background process while our service can still produce to the active (newest) segment and
consume from other segments with no, or at least fewer, conflicts where goroutines access
the same data.

There's always one special segment among the list of segments, and that's the **active segment**.
We call it the active segment because *it's the only segment we actively write to*. When we've
*filled the active segment, we create a new segment and make it the active segment*.

Each segment comprises a store file and an index file. The segment's store file is where we store
the record data; we continually append records to this file. The segment's index file is where we
index each record in the store file.

