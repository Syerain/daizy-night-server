# Project Daizy Night 

Designed for ATOM Reforge Community's daily activities.

The project name Daisy Night is a homophone of The Hazy Night.

It's also my first golang project and I wished to learn something during the work.

## Repo

server & devCliTool
https://github.com/atomreforge/daizy-night-server

scheduled client
https://github.com/atomreforge/daizy-night-app

cli test client
https://github.com/Syerain/dnappcli

## Dev deployment

### Get

```bash
git clone github.com/atomreforge/daizy-night-server
```

### Configuration

during alpha period there are some files for test placed under the root dir.

before first launch, you must :

- edit file `./Config.yaml`

- delete dir `./test`

### Run

```bash
go mod tidy
go run ./cmd/main.go
```
