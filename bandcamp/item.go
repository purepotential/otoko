package bandcamp

import (
	"strconv"
)

type ItemID int64

type ItemType string

const (
	ItemTypeAlbum = "a"
	ItemTypeTrack = "t"
)

func (t ItemType) String() string {
	return string(t)
}

type Item struct {
	// Extraneous metadata stripped, such as fan id,
	// genre, dates, "why", URLs, etc.
	ID        ItemID   `json:"item_id"`
	Type      ItemType `json:"tralbum_type"`
	BandName  string   `json:"band_name"`
	Title     string   `json:"item_title"`
	Purchased Time     `json:"purchased"`
	ArtURL    string   `json:"item_art_url"`
	Sale

	Price    float64 `json:"price"`
	Currency string  `json:"currency"`

	Download string  `json:"-"`
	Tracks   []Track `json:"-"`
}

func (i Item) String() string {
	return string(i.Type) + strconv.FormatInt(int64(i.ID), 10)
}

type SaleID int64

type SaleType string

const (
	// Redeemed from a code
	Code SaleType = "c"

	// Purchased individually from an artist
	Purchase SaleType = "p"

	// Purchased as part of a whole discography
	Records SaleType = "r"
)

type Sale struct {
	ID   SaleID   `json:"sale_item_id"`
	Type SaleType `json:"sale_item_type"`
}

func (s Sale) String() string {
	return string(s.Type) + strconv.FormatInt(int64(s.ID), 10)
}
