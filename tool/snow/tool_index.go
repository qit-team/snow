package main

import "time"

var toolIndexs = []*Tool{
	{
		Name:      "snow",
		Alias:     "snow",
		BuildTime: time.Date(2019, 7, 19, 0, 0, 0, 0, time.Local),
		Install:   "go install github.com/qit-team/snow/tool/snow@latest",
		Summary:   "snow工具集本体",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "snow",
	},
}
