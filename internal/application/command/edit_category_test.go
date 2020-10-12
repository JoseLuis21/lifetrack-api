package command

import (
	"context"
	"fmt"
	"testing"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/category"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
)

func TestNewEditCategoryHandler(t *testing.T) {
	cfg, err := infrastructure.NewConfiguration()
	if err != nil {
		t.Fatal("cannot start configuration")
	}

	r := category.NewInMemoryRepository()

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

	cmd := NewEditCategoryHandler(r, eventbus.NewInMemory(cfg))

	err = cmd.Invoke(EditCategory{
		Ctx:         context.Background(),
		ID:          categories[0].ID,
		Title:       "",
		Description: "",
		Theme:       "",
	})
	if err == nil {
		t.Fatal("edit category command did not failed, expected required field (user)")
	}

	err = cmd.Invoke(EditCategory{
		Ctx:         context.Background(),
		ID:          categories[0].ID,
		Title:       "Special Relativity",
		Description: "Albert Einstein main relativity theory",
		Theme:       "blue",
	})
	if err != nil {
		t.Fatal("edit category command failed", fmt.Sprintf("err: %v", exception.GetDescription(err)))
	}

	t.Log("edit category command succeed")
}
