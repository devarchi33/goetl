package repositories

import (
	"clearance-adapter/config"
)

func getTenantCode() string {
	v := config.GetTenantCode()
	return v
}
