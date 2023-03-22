## g2-sdk-go-base errors

## Error prefixes

Error identifiers are in the format `senzing-PPPPnnnn` where:

`P` is a prefix used to identify the package.
`n` is a location within the package.

Prefixes:

1. `6001` - g2config
1. `6002` - g2configmgr
1. `6003` - g2diagnostic
1. `6004` - g2engine
1. `6005` - g2hasher
1. `6006` - g2product
1. `6007` - g2ssadm

## Common errors

### Postgresql

1. "Error: pq: SSL is not enabled on the server"
    1. The database URL needs the `sslmode` parameter.
       Example:

        ```console
        postgresql://username:password@postgres.example.com:5432/G2/?sslmode=disable
        ```

    1. [Connection String Parameters](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters)
