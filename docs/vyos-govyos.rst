:lastproofread: 2026-03-07

.. _vyos-govyos:

go-vyos
=======

go-vyos is a Go library designed for interacting with VyOS devices through
their REST API. This documentation is intended to guide you in using go-vyos for
programmatic management of your VyOS devices.

- `go-vyos Documentation & Source Code on GitHub <https://github.com/nveeser/go-vyos>`_
  allows you to access and contribute to the library's code.
- `go-vyos on pkg.go.dev <https://pkg.go.dev/github.com/nveeser/go-vyos/vyos>`_ for detailed instructions
  on the installation, configuration, and operation of the go-vyos library.


Installation
------------

You can install go-vyos:

.. code-block:: bash

    go get "github.com/nveeser/go-vyos/vyos"

Getting Started
---------------

Importing and Initializing the Client
-------------------------------------

To initialize a client, provide the host (URL or host:port) and use options for configuration.

.. code-block:: go

    import (
        "context"
        "fmt"
        "github.com/nveeser/go-vyos/vyos"
    )

    client, err := vyos.NewClient("https://192.168.0.1",
        vyos.Token("YOUR_API_KEY"),
        vyos.Insecure(), // Skip TLS verification for self-signed certificates
    )
    if err != nil {
        log.Fatalf("failed to create client: %v", err)
    }

Using Config Mode
-----------------

The Configuration Mode API allows you to modify and retrieve the device configuration.

Set a Configuration
^^^^^^^^^^^^^^^^^^^

.. code-block:: go

    ctx := context.Background()
    err := client.ConfigMode().Set(ctx, "interfaces ethernet eth0 address 192.168.1.1/24")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }

Delete a Configuration
^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: go

    err := client.ConfigMode().Delete(ctx, "interfaces dummy dum1")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }

Show Configuration (Hierarchical)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Returns a map representing the configuration at the given path.

.. code-block:: go

    data, err := client.ConfigMode().Show(ctx, "interfaces ethernet eth0")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("Data: %v\n", data)

Show Values (Multivalue)
^^^^^^^^^^^^^^^^^^^^^^^^

Returns a slice of values for a specific leaf.

.. code-block:: go

    values, err := client.ConfigMode().ShowValues(ctx, "system name-server")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("Values: %v\n", values)

Save and Load Configuration
^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: go

    // Save to default config
    msg, err := client.ConfigMode().Save(ctx, "")
    
    // Save to specific file
    msg, err = client.ConfigMode().Save(ctx, "/config/backup.config")

    // Load from file
    msg, err = client.ConfigMode().Load(ctx, "/config/backup.config")

Using Operation Mode
--------------------

The Operation Mode API allows you to run operational commands.

Show Operational State
^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: go

    output, err := client.OpMode().Show(ctx, "interfaces ethernet eth0")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Println(output)

Generate
^^^^^^^^

.. code-block:: go

    output, err := client.OpMode().Generate(ctx, "pki wireguard key-pair")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Println(output)

Reset
^^^^^

.. code-block:: go

    err := client.OpMode().Reset(ctx, "ip bgp 192.0.2.11")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }

System Info
^^^^^^^^^^^

.. code-block:: go

    info, err := client.OpMode().Info(ctx, vyos.InfoRequest{Version: true, Hostname: true})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    fmt.Printf("Version: %s, Hostname: %s\n", info.Version, info.Hostname)

Advanced Options
----------------

You can configure additional settings like timeouts or custom HTTP clients.

.. code-block:: go

    client, _ := vyos.NewClient("https://vyos",
        vyos.Token("key"),
        vyos.Timeout(10 * time.Second),
        vyos.UserAgent("my-app/1.0"),
        vyos.DebugLogging(), // Prints request/response to stdout
    )
