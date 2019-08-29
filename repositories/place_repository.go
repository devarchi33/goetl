package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
)

// PlaceRepository P2-brand 卖场（场地）仓库
type PlaceRepository struct{}

// GetStoreByCode 根据卖场code获取卖场id
func (PlaceRepository) GetStoreByCode(code string) (models.Store, bool, error) {
	store := new(models.Store)
	has, err := factory.GetP2BrandEngine().Where("tenant_code = ? and code= ? and enable = true", getTenantCode(), code).Get(store)

	return *store, has, err
}
