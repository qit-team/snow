package utils

import (
	"testing"
)

func TestJsonEncode(t *testing.T) {
	//只申明，未初始化
	var p1 map[string]interface{}
	s1, _ := JsonEncode(p1)
	if s1 != "null" {
		t.Error("nil map is not equal {}", s1)
	}

	//已初始化
	p2 := make(map[string]interface{})
	s2, _ := JsonEncode(p2)
	if s2 != "{}" {
		t.Error("blank map is not equal {}", s2)
	}

	//已初始化
	p3 := map[string]interface{}{
		"name": "hts",
	}
	s3, _ := JsonEncode(p3)
	if s3 != "{\"name\":\"hts\"}" {
		t.Error("map is not equal", s3)
	}
}
