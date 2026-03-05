package main

//   Changes
//   * Changed the WithFoo() methods (which return a new client) to use Options func pattern instead.
//      The assumption is that most clients are
//     configured once and reused, and this allows the client to be effectively immutable, removing the need for a mutex.
//   * http.Client is thread-safe so no need to protect it with a mutex.
//
//   * Added the request interface and response values to consolidate the various types of request
//   * Each method uses a request type which can generate an *http.Request
//   * Most use
import (
	"context"
	"github.com/nveeser/go-vyos/vyos"
	"log"
)

func main() {
	ctx := context.Background()
	c, err := vyos.NewClient("https://10.10.10.10",
		vyos.Token("304a441912eb58ebd144d24a584518a46f02123314b209bda06a1315bd8a5425"),
		vyos.Insecure(), vyos.DebugLogging())
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
	{
		resp, err := c.OpMode().Info(ctx, vyos.InfoRequest{})
		if err != nil {
			log.Fatalf("error Info(): %v", err)
		}
		log.Printf("resp: %v", resp)
	}
	{
		resp, err := c.OpMode().Show(ctx, "system image")
		if err != nil {
			log.Fatalf("error Show(): %v", err)
		}
		log.Printf("resp: %v", resp)
	}
	{
		resp, err := c.OpMode().Generate(ctx, "pki wireguard key-pair")
		if err != nil {
			log.Fatalf("error Generate(): %v", err)
		}
		log.Printf("resp: %v", resp)
	}

	{
		resp, err := c.ConfigMode().Show(ctx, "firewall")
		if err != nil {
			log.Fatalf("error Show(): %v", err)
		}
		log.Printf("resp: %v", resp)
	}
	{
		resp, err := c.ConfigMode().ShowValues(ctx, "interfaces ethernet eth1 address")
		if err != nil {
			log.Fatalf("error Show(): %v", err)
		}
		log.Printf("resp: %v", resp)
	}
	{
		resp, err := c.ConfigMode().Exists(ctx, "interfaces ethernet eth1 address")
		if err != nil {
			log.Fatalf("error Exists(): %v", err)
		}
		log.Printf("resp: %v", resp)
	}
	{
		err := c.ConfigMode().Set(ctx, "system option startup-beep")
		if err != nil {
			log.Fatalf("error Set(): %v", err)
		}
		log.Printf("resp: %v", true)
	}
	{
		err := c.ConfigMode().Delete(ctx, "system option startup-beep")
		if err != nil {
			log.Fatalf("error Delete(): %v", err)
		}
		log.Printf("resp: %v", true)
	}
	{
		err := c.ConfigMode().Configure(ctx,
			&vyos.SetRequest{"system option startup-beep"},
			&vyos.DeleteRequest{"system option startup-beep"})
		if err != nil {
			log.Fatalf("error Delete(): %v", err)
		}
		log.Printf("resp: %v", true)
	}
	{
		msg, err := c.ConfigMode().Save(ctx, "")
		if err != nil {
			log.Fatalf("error Save(): %v", err)
		}
		log.Printf("OUT: %v", msg)
	}
}
