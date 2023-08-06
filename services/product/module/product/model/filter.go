package productmodel

type Filter struct {
	ShopId     int    `json:"shop_id,omitempty" form:"shop_id"`
	CategoryId int    `json:"category_id,omitempty" form:"category_id"`
	KeyWord    string `json:"key_word,omitempty" form:"key_word"`
}
