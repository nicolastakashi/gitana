package gitmanager

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	Url           string
	Path          string
	Branch        string
	DashboardPath string
	Auth          RepositoryAuth
}

type RepositoryAuth struct {
	Username string
	Password string
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

	if r.Auth.Username != "" && r.Auth.Password == "" {
		return errors.New("repository.auth.password can not be nil when you inform a username")
	}

	if r.Auth.Password != "" && r.Auth.Username == "" {
		return errors.New("repository.auth.username can not be nil when you inform a password")
	}

	return nil
}

func (r *Repository) Get(ctx context.Context) (bool, error) {
	_, err := git.PlainCloneContext(ctx, r.Path, false, &git.CloneOptions{
		URL:           r.Url,
		ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
		Progress:      os.Stdout,
		Auth:          r.getAuth(),
	})

	if err != nil && err == git.ErrRepositoryAlreadyExists {
		logrus.Info("git repository already exists, trying pull")

		repo, err := git.PlainOpen(r.Path)

		if err != nil {
			logrus.Errorf("error opening git repository folder: %v", err)
			return false, err
		}

		workTree, err := repo.Worktree()

		if err != nil {
			logrus.Errorf("error getting git repository worktree: %v", err)
			return false, err
		}

		err = workTree.PullContext(ctx, &git.PullOptions{
			ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
			Progress:      os.Stdout,
			Auth:          r.getAuth(),
		})

		switch err {
		case nil:
			break
		case git.NoErrAlreadyUpToDate:
			break
		case err.(*os.PathError):

			logrus.Warnf("conflicts with current repo. cloning again. %v", err)

			err = os.RemoveAll(r.Path)

			if err != nil {
				logrus.Errorf("error deleting repo folder: %v", err)
				return false, err
			}

			return r.Get(ctx)
		default:
			logrus.Errorf("error pulling repo: %v", err)
			return false, err
		}
	} else if err != nil {
		logrus.Errorf("error to clone git repository: %v", err)
		return false, err
	}
	return true, nil
}

func (r Repository) getAuth() transport.AuthMethod {
	if r.Auth.Username != "" && r.Auth.Password != "" {
		return &http.BasicAuth{
			Username: r.Auth.Username,
			Password: r.Auth.Password,
		}
	}
	return nil
}

func (r Repository) GetPath() string {
	if r.DashboardPath != "" {
		return filepath.Join(r.Path, r.DashboardPath)
	}
	return r.Path
}