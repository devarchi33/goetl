package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
)

// ProductRepository P2-brand 商品仓库
type ProductRepository struct{}

// GetSkuByCode 根据sku的code获取sku的id
func (ProductRepository) GetSkuByCode(code string) (models.Sku, bool, error) {
	sku := new(models.Sku)
	has, err := factory.GetP2BrandEngine().Where("code= ? and enable = true", code).Get(sku)

	return *sku, has, err
}
