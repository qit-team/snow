package bannerformatter

import (
	"github.com/qit-team/snow/app/models/bannermodel"
)

type BannerFormatter struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Img   string `json:"image"`
	Url   string `json:"url"`
}

func FormatList(bannerList []*bannermodel.Banner) (res []*BannerFormatter) {
	res = make([]*BannerFormatter, len(bannerList))

	for k, banner := range bannerList {
		one := FormatOne(banner)
		res[k] = one
	}

	return res
}

//单条消息的格式化，
func FormatOne(banner *bannermodel.Banner) (res *BannerFormatter) {
	res = &BannerFormatter{
		Id:    int(banner.Id),
		Title: banner.Title,
		Img:   banner.ImageUrl,
		Url:   banner.Url,
	}
	return
}
