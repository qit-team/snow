package entities

// 请求数据结构
type TestRequest struct {
	Name string `json:"name" example:"snow"`
	Url  string `json:"url" example:"github.com/qit-team/snow"`
}

// 返回数据结构
type TestResponse struct {
	Id   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"snow"`
	Url  string `json:"url" example:"github.com/qit-team/snow"`
}

/*
 * validator.v9文档
 * 地址https://godoc.org/gopkg.in/go-playground/validator.v9
 * 列了几个大家可能会用到的，如有遗漏，请看上面文档
 */

// 请求数据结构
type TestValidatorRequest struct {
	// tips，因为组件required不管是没传值或者传 0 or "" 都通过不了，但是如果用指针类型，那么0就是0，而nil无法通过校验
	Id        *int64     `json:"id" validate:"required" example:"1"`
	Age       int        `json:"age" validate:"required,gte=0,lte=130" example:"20"`
	Name      *string    `json:"name" validate:"required" example:"snow"`
	Email     string     `json:"email" validate:"required,email" example:"snow@github.com"`
	Url       string     `json:"url" validate:"required" example:"github.com/qit-team/snow"`
	Mobile    string     `json:"mobile" validate:"required" example:"snow"`
	RangeNum  int        `json:"range_num" validate:"max=10,min=1" example:"3"`
	TestNum   *int       `json:"test_num" validate:"required,oneof=5 7 9" example:"7"`
	Content   *string    `json:"content" example:"snow"`
	Addresses []*Address `json:"addresses" validate:"required,dive,required"`
}

//  Address houses a users address information
type Address struct {
	Street string `json:"street" validate:"required" example:"huandaodonglu"`
	City   string `json:"city" validate:"required" example:"xiamen"`
	Planet string `json:"planet" validate:"required" example:"snow"`
	Phone  string `json:"phone" validate:"required" example:"snow"`
}

// 返回数据结构
type TestValidatorResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}
