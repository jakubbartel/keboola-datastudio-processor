# Keboola Data Studio processor

Convert CSV to a custom Data Studio column based format.

Run tests:
```
go test
```

Build locally:
```
go build -o kbcdatastudioproc cmd/main.go
```

Run locally:
```
KBC_DATADIR=/path/to/your/data/ ./kbcdatastudioproc
```

Build Docker image:
```
docker build -t keboola-datastudio-processor .
```
