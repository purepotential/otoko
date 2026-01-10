package bandcamp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"golang.org/x/net/html"
)

// While bandcamp stores the extension metadata, this is required
// to check if a given track exists for the given format.
var Extensions = map[string]string{
	// Fetched from item download page-data
	"mp3-v0":        ".mp3",
	"mp3-320":       ".mp3",
	"flac":          ".flac",
	"aac-hi":        ".m4a",
	"vorbis":        ".ogg",
	"alac":          ".m4a",
	"wav":           ".wav",
	"aiff-lossless": ".aiff",
}

type Download struct {
	Email    string `json:"-"`
	Size     string `json:"size_mb"`
	Encoding string `json:"encoding_name"`
	URL      string `json:"url"`
}

func (c *Client) GetItemDownload(item *Item, format string) (*Download, error) {
	var data struct {
		// Includes metadata about available encodings, but for simplicity
		// it is stored internally in [Extensions]
		//
		// Bandcamp stores this as an array despite only allowing one
		// item to be downloaded in this request
		Email string `json:"fan_email"`
		Items []struct {
			Downloads map[string]Download `json:"downloads"`
		} `json:"download_items"`
	}

	req, err := http.NewRequest("GET", item.Download, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := decodeBlob(resp.Body, &data); err != nil {
		return nil, err
	}

	digital := &data.Items[0]

	for _, d := range digital.Downloads {
		if d.Encoding != format {
			continue
		}

		d.Email = data.Email
		return &d, nil
	}

	return nil, fmt.Errorf("format %s unavailable", format)

}

func decodeBlob(r io.ReadCloser, v any) error {
	var blob string

	z := html.NewTokenizer(r)
	for {
		t := z.Next()
		if t == html.ErrorToken {
			if z.Err() == io.EOF {
				break
			}

			return z.Err()
		}

		k := z.Token()
		if k.Type != html.StartTagToken {
			continue
		}

		for _, a := range k.Attr {
			if a.Key == "data-blob" {
				blob = a.Val
				break
			}
		}
	}

	if blob == "" {
		return fs.ErrNotExist
	}

	return json.Unmarshal([]byte(blob), v)
}
