package caches

//缓存前缀key，不同的业务使用不同的前缀，避免了业务之间的重用冲突
const (
	Cookie     = "ck:"
	Copy       = "cp:"
	BannerList = "bl:"
)
