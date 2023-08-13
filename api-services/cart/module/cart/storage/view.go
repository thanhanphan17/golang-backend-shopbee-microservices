package cartstorage

import (
	"context"
	"shopbee/common"
)

func (s *cartMySql) ViewMyCart(
	ctx context.Context,
	userId int,
) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	db := s.db

	rows, err := db.Table("carts").
		Select("users.id as shop_id, users.fullname as shop_name, users.avatar as shop_avatar, products.price as price, products.id as product_id").
		Joins("JOIN products ON carts.product_id = products.id").
		Joins("JOIN users ON products.shop_id = users.id").
		Where("carts.user_id = ?", userId).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shopData := make(map[int]map[string]interface{})
	totalPrice := 0
	for rows.Next() {
		var shopID int
		var shopName string
		var shopAvatar *common.Image
		var productID int
		var productPrice int
		err := rows.Scan(&shopID, &shopName, &shopAvatar, &productPrice, &productID)
		if err != nil {
			return nil, err
		}

		if _, ok := shopData[shopID]; !ok {
			totalPrice += productPrice

			shopData[shopID] = map[string]interface{}{
				"shop_id":     shopID,
				"shop_name":   shopName,
				"shop_avatar": shopAvatar,
				"total_price": totalPrice,
				"product_ids": []int{productID},
			}
		} else {
			shopData[shopID]["product_ids"] = append(shopData[shopID]["product_ids"].([]int), productID)
			totalPrice += productPrice
			shopData[shopID]["total_price"] = totalPrice
		}
	}

	for _, data := range shopData {
		result = append(result, data)
	}

	return result, nil
}