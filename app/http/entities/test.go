package entities

//请求数据结构
type TestRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

//返回数据结构
type TestResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}
