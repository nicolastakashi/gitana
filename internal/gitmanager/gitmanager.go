package gitmanager

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	Url           string
	Path          string
	Branch        string
	Proxy         string
	DashboardPath string
	Auth          RepositoryAuth
}

type RepositoryAuth struct {
	Username       string `yaml:"username"`
	AuthSecretName string
	Password       string `yaml:"password"`
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

	if r.Proxy != "" {
		logrus.Debugf("using proxy %v", r.Proxy)

		proxyUrl, err := url.Parse(r.Proxy)

		if err != nil {
			return false, err
		}

		customClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},

			Timeout: 300 * time.Second,
		}

		client.InstallProtocol(proxyUrl.Scheme, githttp.NewClient(customClient))
	}

	gitCloneOptions := &git.CloneOptions{
		URL:           r.Url,
		ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
		Auth:          r.getAuth(),
		Progress:      os.Stdout,
	}

	_, err := git.PlainCloneContext(ctx, r.Path, false, gitCloneOptions)

	if err != nil && err != git.ErrRepositoryAlreadyExists {
		logrus.Error(err)
		return false, err
	}

	if err != nil && err == git.ErrRepositoryAlreadyExists {
		logrus.Debug("git repository already exists, trying pull")

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

		pullOptions := git.PullOptions{
			ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
			Auth:          r.getAuth(),
			Progress:      os.Stdout,
		}

		err = workTree.PullContext(ctx, &pullOptions)

		if err != nil && err != git.NoErrAlreadyUpToDate {
			logrus.Errorf("Could not update the repository %v", err)
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
		return &githttp.BasicAuth{
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
