package command

import (
	"context"
	"fmt"
	"github.com/alexandria-oss/common-go/exception"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/eventbus"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/logging"
	"github.com/neutrinocorp/life-track-api/internal/infrastructure/persistence"
	"testing"
)

func TestNewAddCategoryHandler(t *testing.T) {
	cfg, err := infrastructure.NewConfiguration()
	if err != nil {
		t.Fatal("cannot start configuration")
	}
	logger, cleanup, err := logging.NewZapProd()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	r := persistence.NewCategory(persistence.NewCategoryInMemory(), logger)

	cmd := NewAddCategoryHandler(r, eventbus.NewInMemory(cfg))

	err = cmd.Invoke(AddCategory{
		Ctx:         context.Background(),
		Title:       "",
		User:        "",
		Description: "",
		Theme:       "",
	})
	if err == nil {
		t.Fatal("add category command did not failed, expected required fields (title, user)")
	}

	err = cmd.Invoke(AddCategory{
		Ctx:         context.Background(),
		Title:       "Quantum Mechanics",
		User:        "",
		Description: "",
		Theme:       "",
	})
	if err == nil {
		t.Fatal("add category command did not failed, expected required field (user)")
	}

	err = cmd.Invoke(AddCategory{
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
}
