package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/alecthomas/kong"
	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/chrome"
	_ "github.com/browserutils/kooky/browser/chromium"
	_ "github.com/browserutils/kooky/browser/edge"
	_ "github.com/browserutils/kooky/browser/epiphany"
	_ "github.com/browserutils/kooky/browser/firefox"
	_ "github.com/browserutils/kooky/browser/safari"
	"github.com/sewnie/otoko/bandcamp"
)

type options struct {
	// Pretty unconventional way of setting the client and evaluating the identity cookie.
	Client *Client `kong:"name=identity,help='Bandcamp identity cookie value, fetched from browser if empty',default=,env=BANDCAMP_IDENTITY"`

	Value valueCmd `kong:"cmd,help='Calculate the total value of your Bandcamp collection'"`
	Sync  syncCmd  `kong:"cmd,help='Download and synchronize your collection to a local directory',default=withargs"`
	List  listCmd  `kong:"cmd,help='Display detailed metadata for tracks and albums in your collection'"`
}

func main() {
	var o options
	app := kong.Parse(&o,
		kong.UsageOnError())

	err := app.Run(o.Client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (o *options) AfterApply() error {
	f, err := o.Client.GetFan()
	if err != nil {
		return err
	}
	o.Client.Fan = f
	return nil
}

type Client struct {
	Fan *bandcamp.Fan

	*bandcamp.Client
}

func (c *Client) UnmarshalText(b []byte) error {
	identity := string(b)

	if identity == "" {
		ctx := context.Background()
		cookies := kooky.TraverseCookies(ctx,
			kooky.Valid,
			kooky.DomainHasSuffix(`bandcamp.com`),
			kooky.FilterFunc(func(c *kooky.Cookie) bool {
				return c.Name == "identity"
			}),
		).Collect(ctx)
		sort.SliceStable(cookies, func(i, j int) bool {
			return cookies[i].Expires.After(cookies[j].Expires)
		})
		if len(cookies) > 0 {
			identity = cookies[0].Value
		}
	}
	if identity == "" {
		return errors.New("bandcamp identity required")
	}

	*c = Client{Client: bandcamp.New(identity)}
	return nil
}
