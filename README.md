# Gallery Service

This is a simple gallery service written in golang and run in docker.

## Usage

### Start MariaDB for local development

```
make mariadb
```

### Initialize DB

```
make initdb
```

### Start Service

```
make up
```

### Check Logs

```
make logs
```

### Stop Service

```
make down
```

## TODO

- replace golang template by some javaScript framework like react or vuejs
- add unit tests
