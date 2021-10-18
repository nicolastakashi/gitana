package dashboardloader

import (
	"testing"
)

func TestLoadInvalidPath(t *testing.T) {
	dashboards, err := Load("./dev/blabla")

	if err == nil {
		t.Fatal("load an path file must return an error")
	}

	if len(dashboards) > 0 {
		t.Fatal("load an invalid path must not return a dashboard struct")
	}
}

func TestLoadInvalidFile(t *testing.T) {
	dashboards, err := Load("./testdata/invalid")

	if err == nil {
		t.Fatal("load an invalid json file must return an error")
	}

	if len(dashboards) > 0 {
		t.Fatal("load an invalid json file must not return a dashboard struct")
	}
}

func TestLoadValidFile(t *testing.T) {
	dashboards, err := Load("./testdata/valid")

	if err != nil {
		t.Fatal("load an valid json file must not return an error")
	}

	if len(dashboards) == 0 {
		t.Fatal("load an valid json file must return a dashboard struct")
	}
}

func TestLoadEmptyFile(t *testing.T) {
	dashboards, err := Load("./testdata/empty")

	if err == nil {
		t.Fatal("load an empty json file must return an error")
	}

	if len(dashboards) > 0 {
		t.Fatal("load an empty json file must not return a dashboard struct")
	}
}
