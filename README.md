# go-vyos

go-vyos is a Go client library for accessing the [VyOS REST API](https://docs.vyos.io/en/latest/automation/vyos-api.html#vyosapi)

## Usage

```go
import "github.com/nveeser/go-vyos/vyos"
```

Construct a new VyOS client by providing the host URL (or host:port) and options like authentication token and TLS configuration.

```go
c, err := vyos.NewClient("https://192.168.0.1", 
    vyos.Token("AUTH_KEY"),
    vyos.Insecure(), // Use for self-signed certificates
)
if err != nil {
    log.Fatal(err)
}
```

The client separates functionality into `ConfigMode()` and `OpMode()`.

### Configuration Mode (ConfigMode)

#### Set a Configuration Value

```go
err := c.ConfigMode().Set(ctx, "interfaces ethernet eth0 address 192.168.1.1/24")
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

#### Delete a Configuration Object

```go
err := c.ConfigMode().Delete(ctx, "interfaces dummy dum1")
```

#### Show Configuration Data

```go
// Show returns a map[string]any representation of the configuration
data, err := c.ConfigMode().Show(ctx, "interfaces dummy dum1")
if err == nil {
    fmt.Printf("Data: %v\n", data)
}

// ShowValues returns a slice of values
values, err := c.ConfigMode().ShowValues(ctx, "interfaces ethernet eth0 address")
```

#### Batch Configuration

```go
err := c.ConfigMode().Configure(ctx, 
    &vyos.SetRequest{Path: "interfaces dummy dum1"},
    &vyos.SetRequest{Path: "interfaces dummy dum1 address 10.0.0.1/24"},
)
```

#### Save and Load

```go
// Save current configuration to the default location
msg, err := c.ConfigMode().Save(ctx, "")

// Save to a specific file
msg, err := c.ConfigMode().Save(ctx, "/config/backup.config")

// Load from a file
msg, err := c.ConfigMode().Load(ctx, "/config/backup.config")
```

### Operational Mode (OpMode)

#### Show Operational Data

```go
output, err := c.OpMode().Show(ctx, "system image")
if err == nil {
    fmt.Println(output)
}
```

#### Generate Data

```go
output, err := c.OpMode().Generate(ctx, "pki wireguard key-pair")
```

#### Reset a Service

```go
err := c.OpMode().Reset(ctx, "ip bgp 192.0.2.11")
```

#### System Information

```go
info, err := c.OpMode().Info(ctx, vyos.InfoRequest{Version: true, Hostname: true})
if err == nil {
    fmt.Printf("Version: %s, Hostname: %s\n", info.Version, info.Hostname)
}
```

## License

This library is distributed under the BSD-style license found in the [LICENSE](./LICENSE)
file.
