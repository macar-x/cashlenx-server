package model

type CategoryDTO struct {
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
}
