package trees

import (
	"tree_manager/internal/storage"

	"github.com/rs/zerolog"
)

type TreeManager struct {
	Repo   *storage.Repo
	Logger *zerolog.Logger
}

func NewTreeManager(repo *storage.Repo, logger *zerolog.Logger) *TreeManager {
	return &TreeManager{Repo: repo, Logger: logger}
}
