package gitmanager

import (
	"context"
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	Url    string
	Path   string
	Branch string
}

func (r Repository) Validate() error {
	if r.Url == "" {
		return errors.New("repository.url flag can not be null")
	}

	if r.Branch == "" {
		return errors.New("repository.branch flag can not be null")
	}

	if r.Path == "" {
		return errors.New("repository.branch flag can not be null")
	}
	return nil
}

func (gm *Repository) Get(ctx context.Context) (bool, error) {
	_, err := git.PlainCloneContext(ctx, gm.Path, false, &git.CloneOptions{
		URL:           gm.Url,
		ReferenceName: plumbing.NewBranchReferenceName(gm.Branch),
		Progress:      os.Stdout,
	})

	if err != nil && err == git.ErrRepositoryAlreadyExists {
		logrus.Info("git repository already exists, trying pull")

		repo, err := git.PlainOpen(gm.Path)

		if err != nil {
			logrus.Error("error opening git repository folder: %v", err)
			return false, err
		}

		workTree, err := repo.Worktree()

		if err != nil {
			logrus.Error("error getting git repository worktree: %v", err)
			return false, err
		}

		err = workTree.PullContext(ctx, &git.PullOptions{
			ReferenceName: plumbing.NewBranchReferenceName(gm.Branch),
			Progress:      os.Stdout,
		})

		switch err {
		case nil:
			break
		case git.NoErrAlreadyUpToDate:
			break
		case err.(*os.PathError):

			logrus.Warn("conflicts with current repo. cloning again. %v", err)

			err = os.RemoveAll(gm.Path)

			if err != nil {
				logrus.Error("error deleting repo folder: %v", err)
				return false, err
			}

			return gm.Get(ctx)
		default:
			logrus.Error("error pulling repo: %v", err)
			return false, err
		}
	} else if err != nil {
		logrus.Error("error to clone git repository: %v", err)
		return false, err
	}
	return true, nil
}
