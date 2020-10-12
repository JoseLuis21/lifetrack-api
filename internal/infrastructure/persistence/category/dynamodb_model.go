package category

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"
)

// DynamoModel AWS DynamoDB model
type DynamoModel struct {
	PK          string `json:"pk"`
	SK          string `json:"sk"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	User        string `json:"user"`
	Color       string `json:"color,omitempty"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	Active      bool   `json:"active"`
	GSIPK       string `json:"gsipk"`
	GSISK       string `json:"gsisk"`
}

func (c DynamoModel) ToModel() *model.Category {
	return &model.Category{
		ID:          util.FromDynamoID(c.PK),
		Title:       c.Title,
		Description: c.Description,
		User:        util.FromDynamoID(c.GSIPK),
		Color:       c.Color,
		CreateTime:  c.CreateTime,
		UpdateTime:  c.UpdateTime,
		Active:      c.Active,
	}
}
