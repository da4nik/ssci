package store

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"

	"github.com/da4nik/ssci/config"
	"github.com/da4nik/ssci/types"
)

const bucketName = "projects"
const dbFileName = "ssci.db"

var dbFilePath = filepath.Join(config.Workspace, dbFileName)

// LoadProject - loads project
func LoadProject(name string) (*types.Project, error) {
	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var project types.Project
	err = db.Update(func(tx *bolt.Tx) error {
		b, terr := tx.CreateBucketIfNotExists([]byte(bucketName))
		if terr != nil {
			return terr
		}

		v := b.Get([]byte(name))
		if v == nil {
			return fmt.Errorf("Project %s not found", name)
		}

		if unmarshalErr := json.Unmarshal(v, &project); unmarshalErr != nil {
			return unmarshalErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// SaveProject saves build to local storage
func SaveProject(project *types.Project) error {
	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, terr := tx.CreateBucketIfNotExists([]byte(bucketName))
		if terr != nil {
			return terr
		}

		json, marshalErr := json.Marshal(project)
		if marshalErr != nil {
			return marshalErr
		}

		return b.Put([]byte(project.Name), json)
	})
	if err != nil {
		return err
	}
	return nil
}

// NewBuild return new build for project
func NewBuild(project *types.Project) *types.Build {
	build := types.Build{
		ID:        nextBuildID(project),
		StartTime: time.Now(),
	}
	project.Builds = append(project.Builds, build)
	return &build
}

// NewProject return new project entity
func NewProject(name, url string) *types.Project {
	return &types.Project{
		Name: name,
		Repo: url,
	}
}

func nextBuildID(project *types.Project) int {
	if len(project.Builds) == 0 {
		return 0
	}

	buildNumber := 0
	for _, build := range project.Builds {
		if build.ID > buildNumber {
			buildNumber = build.ID
		}
	}
	return buildNumber + 1
}

func log() *logrus.Entry {
	return logrus.WithField("module", "store")
}
