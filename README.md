# blob-retriever
restore pruned blob for prysm

## Usage

```
NAME:
   blob_retriever - Retrieve and check pruned blobs

USAGE:
   blob_retriever [options] command

COMMANDS:
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --mode value, -m value       run mode (retrieve / check)
   --api_url value, -u value    Beacon node URL
   --api_type value, -a value   Beacon node network type (any or prysm)
   --data_path value, -d value  data path to store blobs
   --worker value, -w value     number of workers
   --from value, -f value       from slot. minimum is 8626176
   --to value, -t value         to slot
   --help, -h                   show help
```

## Build and run

    make all