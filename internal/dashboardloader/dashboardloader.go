package dashboardloader

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Dashboard struct {
	Name      string
	FileName  string
	Folder    string
	Dashboard map[string]interface{}
}

func Load(path string) map[string]Dashboard {
	dashboards := map[string]Dashboard{}

	err := filepath.Walk(path, func(path string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() && fi.Name() == ".git" {
			return filepath.SkipDir
		}

		if !fi.IsDir() && filepath.Ext(path) == ".json" {
			logrus.Debugf("start to load dashboard %v", fi.Name())

			dash := Dashboard{
				FileName: fi.Name(),
				Name:     getSanitizedName(fi.Name()),
				Folder:   getDashboardDirName(path),
			}

			err := readDashboardFile(path, &dash.Dashboard)

			if err != nil {
				logrus.Errorf("error to load dashboard file %v error: %v", dash.Name, err)
				return err
			}

			logrus.Debugf("dashboard %v content readed", dash.Name)

			dashboards[dash.Name] = dash
		}

		return nil
	})

	if err != nil {
		return nil
	}

	return dashboards
}

func readDashboardFile(path string, value *map[string]interface{}) error {

	dashboardFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer dashboardFile.Close()

	byteValue, _ := ioutil.ReadAll(dashboardFile)

	json.Unmarshal([]byte(byteValue), &value)

	return nil
}

func getDashboardDirName(path string) string {
	dirs := strings.Split(filepath.Dir(path), "/")
	return dirs[len(dirs)-1]
}

func getSanitizedName(name string) string {
	m := regexp.MustCompile("[^a-zA-Z0-9]+|json")
	return m.ReplaceAllString(name, "")
}

func (d *Dashboard) ToConfigMap(namespace string, dashBoardLabels string, folderAnnotation string) (v1.ConfigMap, error) {

	json, err := json.MarshalIndent(d.Dashboard, "", " ")

	if err != nil {
		logrus.Error("error to marshal dashboard content %v", err)
		return v1.ConfigMap{}, err
	}

	cm := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: namespace,
			Labels:    parseDashboardLabels(dashBoardLabels),
		},
		Data: map[string]string{
			d.FileName: string(json),
		},
	}

	if folderAnnotation != "" {
		cm.Annotations = map[string]string{
			folderAnnotation: d.Folder,
		}
	}

	return cm, nil
}

func parseDashboardLabels(dashBoardLabels string) map[string]string {

	keyValueLabels := strings.Split(dashBoardLabels, "=")

	labels := map[string]string{
		"app.kubernetes.io/managed-by": "gitana",
	}

	for i := 0; i < len(keyValueLabels); i += 2 {
		labels[keyValueLabels[i]] = keyValueLabels[i+1]
	}

	return labels
}

func (d *Dashboard) NeedsToUpdate(ccm *v1.ConfigMap, dcm *v1.ConfigMap) bool {

	if !reflect.DeepEqual(ccm.ObjectMeta.Labels, dcm.ObjectMeta.Labels) {
		return true
	}

	if !reflect.DeepEqual(ccm.ObjectMeta.Annotations, dcm.ObjectMeta.Annotations) {
		return true
	}

	if !reflect.DeepEqual(ccm.Data, dcm.Data) {
		return true
	}

	return false
}
