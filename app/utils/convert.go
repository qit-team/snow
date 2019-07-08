package utils

import (
	"fmt"
)

func SliceStr2Interface(input []string) (output []interface{}) {
	for _, v := range input {
		output = append(output, v)
	}
	return
}

func MapStrInterface2MapStrStr(input map[string]interface{}) (output map[string]string) {
	output = make(map[string]string)
	for k, v := range input {
		output[k] = Interface2Str(v)
	}
	return
}

//interface转换为字符串
func Interface2Str(v interface{}) string {
	switch v.(type) {
	case float64, float32:
		return fmt.Sprintf("%f", v)
	case []rune:
		return string(v.([]rune))
	default:
		return fmt.Sprint(v)
	}
}

func Num2Str(u interface{}) string {
	return fmt.Sprintf("%d", u)
}
