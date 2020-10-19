package category

import (
	"time"

	"github.com/neutrinocorp/life-track-api/internal/domain/aggregate"
	"github.com/neutrinocorp/life-track-api/internal/domain/entity"
	"github.com/neutrinocorp/life-track-api/internal/domain/model"
	"github.com/neutrinocorp/life-track-api/internal/domain/value"
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
	Image       string `json:"image,omitempty"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	Active      bool   `json:"active"`
	GSIPK       string `json:"gsipk"`
	GSISK       string `json:"gsisk"`
}

func (c DynamoModel) ToAggregate() *aggregate.Category {
	ag := new(aggregate.Category)
	_ = ag.AssignUser(c.User)

	id := new(value.CUID)
	_ = id.Set(c.PK)

	title := new(value.Title)
	_ = title.Set(c.Title)

	description := new(value.Description)
	_ = description.Set(c.Description)

	color := new(value.Color)
	_ = color.Set(c.Color)

	image := new(value.Image)
	_ = color.Set(c.Image)

	metadata := new(value.Metadata)
	metadata.SetCreateTime(time.Unix(c.CreateTime, 0).UTC())
	metadata.SetUpdateTime(time.Unix(c.UpdateTime, 0).UTC())
	metadata.SetState(c.Active)

	ag.Set(&entity.Category{
		ID:          id,
		Title:       title,
		Description: description,
		Color:       color,
		Image:       image,
		Metadata:    metadata,
	})

	return ag
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
