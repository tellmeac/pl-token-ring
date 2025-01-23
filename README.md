# Token ring

Лабораторная работа.

## Launch

```bash
go run . 500
```

## Example

```bash
go run ./cmd/server -n 12
```

```bash
curl -d '{"data":"hello world", "ttl": 100, "receiver": 0}' http:/localhost:8080
```
