# blob-retriever
restore pruned blob for prysm

## Usage

```
NAME:
   blob_retriever - Retrieve and check blobs

USAGE:
   blob_retriever [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --mode value, -m value         run mode (retrieve / check) (default: "retrieve")
   --beacon_url value, -b value   Beacon node URL
   --beacon_type value, -n value  Beacon node network type (any or prysm) (default: "any")
   --data value, -d value         data path to store blobs
   --worker value, -w value       number of workers (default: 1)
   --from value, -f value         from slot (default: 0)
   --to value, -t value           to slot (default: 0)
   --help, -h                     show help
```

## Build and run

    make all