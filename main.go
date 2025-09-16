package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/chrome"
	_ "github.com/browserutils/kooky/browser/chromium"
	_ "github.com/browserutils/kooky/browser/edge"
	_ "github.com/browserutils/kooky/browser/epiphany"
	_ "github.com/browserutils/kooky/browser/firefox"
	_ "github.com/browserutils/kooky/browser/safari"
)

var (
	dir         string
	jobs        int64
	format      string
	cookiesFile string
)

func init() {
	flag.StringVar(&format, "format", "flac", "audio format to use")
	flag.Int64Var(&jobs, "jobs", 6, "amount of parallel jobs to use to download")
	flag.StringVar(&dir, "o", "collection", "directory to download albums to")
	flag.StringVar(&cookiesFile, "cookies", "otoko-cookies.txt", "bandcamp user cookies file path")
}

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cookie, err := getCookie(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(22)
	}

	downloader, err := new(ctx, cookie)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = downloader.downloadCollection()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(11)
	}
}

// TODO: turn into []http.Cookie or http.CookieJar
func getCookie(ctx context.Context) (string, error) {
	cookies := kooky.TraverseCookies(ctx,
		kooky.Valid,
		kooky.DomainHasSuffix(`bandcamp.com`),
		kooky.FilterFunc(func(c *kooky.Cookie) bool {
			return c.Name == "identity" || c.Name == "session"
		}),
	).Collect(ctx)
	if len(cookies) > 0 {
		var s []string
		for _, c := range cookies {
			// kooky adds optional values to the final http.Cookie
			// which makes its Stringer return for Set-Cookie
			s = append(s, c.Name+"="+c.Value)
		}
		return strings.Join(s, "; "), nil
	}
	if cookiesFile == "" {
		return "", errors.New("cookie file missing")
	}

	v, err := os.ReadFile(cookiesFile)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(v)), nil
}
