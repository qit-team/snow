package bannerservice

import (
	"github.com/qit-team/snow/app/models/bannermodel"
)

func GetListByPid(pid int, limit int, page int) (banners []*bannermodel.Banner, err error) {
	limitStart := GetLimitStart(limit, page)
	banners, err = bannermodel.GetInstance().GetListByPid(pid, limitStart...)
	return
}

func GetLimitStart(limit int, page int) (arr []int) {
	arr = make([]int, 2)
	if limit <= 0 {
		limit = 20
	}
	arr[0] = limit
	if page > 0 {
		arr[1] = (page - 1) * limit
	} else {
		arr[1] = 0
	}
	return
}
