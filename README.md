# blob-retriever
restore pruned blob for prysm

## Usage

```
Usage:
  blob_retriever [flags]

Flags:
  -n, --beacon_type string   Beacon node network type. ( any or prysm ) (default "any")
  -b, --beacon_url string    Beacon node URL
  -d, --data string          data path to store blobs
  -f, --from uint            start slot
  -h, --help                 help for blob_retriever
  -m, --mode string          run mode (retrieve / check) (default "retrieve")
  -t, --to uint              end slot
  -w, --worker uint8         number of worker (default 1)
```