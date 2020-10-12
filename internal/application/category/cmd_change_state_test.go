package category

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence/category"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/domain/shared"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
)

func TestNewChangeCategoryStateHandler(t *testing.T) {
	cfg, err := infrastructure.NewConfiguration()
	if err != nil {
		t.Fatal("cannot start configuration")
	}

	r := category.NewInMemoryRepository()
	b := eventbus.NewInMemory(cfg)

	cmdAdd := NewAddHandler(r, b)
	err = cmdAdd.Invoke(Add{
		Ctx:         context.Background(),
		Title:       "Classical Mechanics",
		User:        "123",
		Description: "",
		Theme:       "red",
	})
	if err != nil {
		t.Fatal("add category command failed", fmt.Sprintf("err: %v", exception.GetDescription(err)))
	}

	t.Log("add category command succeed")

	categories, _, err := r.Fetch(context.Background(), "", 1, shared.CategoryCriteria{})
	if err != nil {
		t.Fatal("list category query failed", fmt.Sprintf("err: %v", exception.GetDescription(err)))
	}

	t.Log("list category query succeed")
	t.Log(categories[0])

	cmd := NewChangeStateHandler(r, b)
	err = cmd.Invoke(ChangeState{
		Ctx:   context.Background(),
		ID:    "",
		State: false,
	})
	if err == nil || errors.Is(err, exception.RequiredField) {
		t.Fatal("change state category command failed, expected required field (id)")
	}

	err = cmd.Invoke(ChangeState{
		Ctx:   context.Background(),
		ID:    categories[0].ID,
		State: false,
	})
	if err != nil {
		t.Fatal("change state category command failed, expected nil error", fmt.Sprintf("err: %v", exception.GetDescription(err)))
	}

	t.Log("change state category command succeed")
}
