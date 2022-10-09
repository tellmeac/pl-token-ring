# Token Ring. 

Laboratory work

## Launch

```bash
export ELEMENTS_COUNT=100
go run . $ELEMENTS_COUNT
```

## Posting tokens 

```bash
curl -d '{"data":"hello world", "ttl": 100, "reciever": 0}' http:/localhost:8080
```