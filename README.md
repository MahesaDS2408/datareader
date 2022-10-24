## Datareader

An tool to convert CSV data to customable read-only API.

## API

| URL | Desc | Example
|--|--|--|
| `/<csv filename>` | This will show all csv file list | `/test` |
| `/<csv filename>?allowed=<>` | This will show filtered table csv file list | `/test?allowed=nama` |

## How to Use

### To run on local

- `go run main.go`

### Deploy to server

- `go build .` on you local computer.
- Copy the binary `datareader` and `.yodelconf.tf` to server.
- Config the `.yodelconf.tf`, the default you must create folder named `test` and put all csv you need there.

### Maintainer

- [@frederett](https://github.com/frederett), [Donate to me](https://github.com/sponsors/frederett).