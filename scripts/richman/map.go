package richman

type MapElement struct {
	ID          string         `json:"id"`          // ID
	Name        string         `json:"name"`        // 名称
	LocationX   int            `json:"location_x"`  // 土地X坐标
	LocationY   int            `json:"location_y"`  // 土地位置Y坐标
	Level       int            `json:"level"`       // 土地星级
	Fee         int64          `json:"fee"`         // 购买基础费用
	RentFee     int64          `json:"rent_fee"`    // 租金
	Enable      int            `json:"enable"`      // 标记是否可用，已购买，空地，已被抵押状态下不可用
	Flag        GAME_FLAG_ENUM `json:"flag"`        // 标记
	Description string         `json:"description"` // 简介
}

func (m MapElement) IsEqual(land MapElement) bool {

	if m.LocationX == land.LocationX && m.LocationY == land.LocationY {
		return true
	}
	return false
}
