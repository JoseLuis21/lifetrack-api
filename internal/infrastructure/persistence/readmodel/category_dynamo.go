package readmodel

import (
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/util"
)

// CategoryDynamo read model
type CategoryDynamo struct {
	PK          string `json:"pk"`
	SK          string `json:"sk"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	User        string `json:"user"`
	Theme       string `json:"theme,omitempty"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	Active      bool   `json:"active"`
	GSIPK       string `json:"gsipk"`
	GSISK       string `json:"gsisk"`
}

func NewCategoryDynamo(schemaName string, m model.Category) *CategoryDynamo {
	return &CategoryDynamo{
		PK:          util.GenerateDynamoID(schemaName, m.ID),
		SK:          util.GenerateDynamoID(schemaName, m.ID),
		Title:       m.Title,
		Description: m.Description,
		User:        m.User,
		Theme:       m.Theme,
		CreateTime:  m.CreateTime,
		UpdateTime:  m.UpdateTime,
		Active:      m.Active,
		GSIPK:       util.GenerateDynamoID("User", m.User),
		GSISK:       util.GenerateDynamoID(schemaName, m.ID),
	}
}

func (c CategoryDynamo) ToModel() *model.Category {
	return &model.Category{
		ID:          util.FromDynamoID(c.PK),
		Title:       c.Title,
		Description: c.Description,
		User:        util.FromDynamoID(c.GSIPK),
		Theme:       c.Theme,
		CreateTime:  c.CreateTime,
		UpdateTime:  c.UpdateTime,
		Active:      c.Active,
	}
}
