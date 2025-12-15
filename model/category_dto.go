package model

type CategoryDTO struct {
	ParentName string `json:"parent_name"`
	Name       string `json:"name"`
	Remark     string `json:"remark"`
}
