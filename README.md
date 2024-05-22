# blob_restorer
restore pruned blob for prysm and lighthouse

## Usage

```
Usage:
  blob_restorer [flags]

Flags:
  -b, --beacon string         Beacon node URL
  -f, --from uint             Start slot
  -h, --help                  help for blob_retriever
  -m, --mode string           run mode (retrieve / check) (default "retrieve")
  -p, --storage_path string   Path to store blobs
  -t, --to uint               End slot
  -w, --worker uint8          Number of worker (default 1)
```