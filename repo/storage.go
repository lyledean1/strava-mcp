package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"stravamcp/model"
	"strings"
)

type Storage interface {
	GetAllAthleteActivities() ([]model.AthleteActivity, error)
	GetAthleteActivity(id string) (*model.AthleteActivity, error)
	GetAllActivityStreams() ([]model.ActivityStreams, error)
	GetActivityStream(id string) (*model.ActivityStreams, error)
	SaveAthleteActivity(activity *model.AthleteActivity) error
	SaveActivityStream(id string, stream *model.ActivityStreams) error
}
type storage struct {
	path string
}

func NewStorage(path string) Storage {
	return &storage{path: path}
}

func (s *storage) GetAthleteActivity(id string) (*model.AthleteActivity, error) {
	var loadedData model.AthleteActivity
	err := LoadFromZstd(s.getFilePath(id, "activity"), &loadedData)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &loadedData, nil
}

func (s *storage) GetAllAthleteActivities() ([]model.AthleteActivity, error) {
	dirPath := fmt.Sprintf("%s/data/%s/", s.path, "activity")

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var activities []model.AthleteActivity

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".zstd") {
			continue
		}

		fullPath := filepath.Join(dirPath, fileName)

		var activity model.AthleteActivity
		if err := LoadFromZstd(fullPath, &activity); err != nil {
			continue
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (s *storage) GetAllActivityStreams() ([]model.ActivityStreams, error) {
	//TODO implement me
	panic("implement me")
}

func (s *storage) GetActivityStream(id string) (*model.ActivityStreams, error) {
	var loadedData model.ActivityStreams
	err := LoadFromZstd(s.getFilePath(id, "stream"), &loadedData)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &loadedData, nil
}

func (s *storage) SaveAthleteActivity(activity *model.AthleteActivity) error {
	return SaveToZstd(activity, s.getFilePath(fmt.Sprintf("%d", activity.ID), "activity"))
}

func (s *storage) SaveActivityStream(id string, stream *model.ActivityStreams) error {
	return SaveToZstd(stream, s.getFilePath(id, "stream"))
}

func (s *storage) getFilePath(id, activityType string) string {
	return fmt.Sprintf("%s/data/%s/%s.json.zstd", s.path, activityType, id)
}
