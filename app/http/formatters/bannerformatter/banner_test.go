package bannerformatter

import (
	"github.com/qit-team/snow/app/models/bannermodel"
	"testing"
)

func TesFormatOne(t *testing.T) {
	a := &bannermodel.Banner{
		Id:       1,
		Title:    "test",
		ImageUrl: "http://x/1.jpg",
		Url:      "http://x",
		Status:   "1",
	}
	b := FormatOne(a)
	if b.Title != a.Title || b.Img != a.ImageUrl || b.Url != a.Url {
		t.Error("FormatOne not same")
	}
}

func TesFormatList(t *testing.T) {
	a := make([]*bannermodel.Banner, 2)
	a[0] = &bannermodel.Banner{
		Id:       1,
		Title:    "test",
		ImageUrl: "http://x1/1.jpg",
		Url:      "http://x1",
		Status:   "1",
	}
	a[1] = &bannermodel.Banner{
		Id:       2,
		Title:    "test2",
		ImageUrl: "http://x/2.jpg",
		Url:      "http://x2",
		Status:   "2",
	}
	b := FormatList(a)
	for k, v := range b {
		if v.Title != a[k].Title || v.Img != a[k].ImageUrl || v.Url != a[k].Url {
			t.Error("FormatList not same")
		}
	}
}
