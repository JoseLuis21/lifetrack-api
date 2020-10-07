package command

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
)

func TestNewRemoveCategoryHandler(t *testing.T) {
	cfg, err := infrastructure.NewConfiguration()
	if err != nil {
		t.Fatal("cannot start configuration")
	}

	r := persistence.NewCategoryInMemory()

	cmdAdd := NewAddCategoryHandler(r, eventbus.NewInMemory(cfg))
	err = cmdAdd.Invoke(AddCategory{
		Ctx:         context.Background(),
		Title:       "Classical Mechanics",
		User:        "123",
		Description: "",
		Theme:       "red",
	})
	if err != nil {
		t.Fatal("add category command failed", fmt.Sprintf("err: %v", exception.GetDescription(err)))
	}

	categories, _, err := r.Fetch(context.Background(), "", 1, shared.CategoryCriteria{})
	if err != nil {
		t.Fatal("list category query failed", fmt.Sprintf("err: %v", exception.GetDescription(err)))
	}

	t.Log("list category query succeed")
	t.Log(categories[0])

	cmd := NewRemoveCategoryHandler(r, eventbus.NewInMemory(cfg))
	err = cmd.Invoke(RemoveCategory{
		Ctx: context.Background(),
		ID:  "",
	})
	if err == nil {
		t.Fatal("remove category command did not failed, expected required field (id)")
	}

	err = cmd.Invoke(RemoveCategory{
		Ctx: context.Background(),
		ID:  categories[0].ID,
	})
	if err != nil {
		t.Fatal("remove category command failed", fmt.Sprintf("err: %v", exception.GetDescription(err)))
	}

	categories, _, err = r.Fetch(context.Background(), "", 1, shared.CategoryCriteria{})
	if err == nil || !errors.Is(err, exception.NotFound) {
		t.Fatal("remove category command failed, expected category not found")
	}

	t.Log("remove category command succeed")
}
