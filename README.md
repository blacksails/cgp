# CGP

This is a go api wrapper for the excellent Communigate Pro mailserver.

## Usage

Create an api client

```go
c := cgp.New("https://server.hostname/CLI/", "username", "password")
```

For API instuctions of the client see the go-docs

## Contributing
Not all API calls are wrapped yet. I have added the ones that I needed more
will be added if need be. Please submit an issue or a pull request.
