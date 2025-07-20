package service

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"stravamcp/model"
	"stravamcp/pkg/client"
	"stravamcp/repo"
	"strings"
	"time"
)

type ActivityService interface {
	ProcessActivities(after time.Time) error
	GetAllActivities(_ context.Context, filter string, before *time.Time, after *time.Time) ([]model.AthleteActivity, error)
	GetActivityStream(_ context.Context, id string) (*ActivityStreamData, error)
}
type activityService struct {
	stravaClient client.StravaClient
	tokenRepo    repo.TokenRepo
	storage      repo.Storage
}

func NewActivityService(stravaClient client.StravaClient, tokenRepo repo.TokenRepo, storage repo.Storage) ActivityService {
	return &activityService{stravaClient: stravaClient, tokenRepo: tokenRepo, storage: storage}
}

func (a *activityService) ProcessActivities(after time.Time) error {
	token, err := a.tokenRepo.Get()
	if err != nil {
		return err
	}
	activities, err := a.stravaClient.GetAllAthleteActivities(int(after.Unix()), token.AccessToken)
	if err != nil {
		return err
	}

	for _, athleteActivity := range activities {
		id := athleteActivity.ID
		activity, err := a.storage.GetAthleteActivity(fmt.Sprintf("%d", id))
		if err != nil {
			return err
		}
		if activity == nil {
			err = a.storage.SaveAthleteActivity(&athleteActivity)
			if err != nil {
				return err
			}
		}
		activityStream, err := a.storage.GetActivityStream(fmt.Sprintf("%d", id))
		if err != nil {
			return err
		}
		if activityStream != nil {
			continue
		}
		slog.Info("Getting stream for activity", "id", id, "start_date", athleteActivity.StartDate)
		stream, err := a.stravaClient.FetchStreams(fmt.Sprintf("%d", id), getActivityKeys(), token.AccessToken)
		if err != nil {
			return err
		}
		err = a.storage.SaveActivityStream(fmt.Sprintf("%d", id), stream)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *activityService) GetAllActivities(_ context.Context, filter string, before *time.Time, after *time.Time) ([]model.AthleteActivity, error) {
	defaultAfter := after
	if defaultAfter == nil {
		oneDayAgo := time.Now().Add(time.Hour * -24 * 7)
		defaultAfter = &oneDayAgo
	}
	err := a.ProcessActivities(*defaultAfter)
	if err != nil {
		return nil, err
	}
	allActivities, err := a.storage.GetAllAthleteActivities()
	if err != nil {
		return nil, err
	}

	slices.SortFunc(allActivities, func(a, b model.AthleteActivity) int {
		return cmp.Compare(b.StartDate, a.StartDate)
	})

	filteredActivities := make([]model.AthleteActivity, 0)

	for _, activity := range allActivities {
		if filter != "" && strings.EqualFold(activity.Type, filter) {
			continue
		}

		activityDate, err := time.Parse(time.RFC3339, activity.StartDate)
		if err != nil {
			continue
		}

		if after != nil && activityDate.Before(*after) {
			continue
		}

		if before != nil && activityDate.After(*before) {
			continue
		}

		filteredActivities = append(filteredActivities, activity)
	}

	return filteredActivities, nil
}

type ActivityStreamData struct {
	ActivityID   string            `json:"activity_id"`
	ActivityName string            `json:"activity_name,omitempty"`
	Date         string            `json:"date,omitempty"`
	Streams      []StreamDataPoint `json:"streams"`
}

type StreamDataPoint struct {
	Time      *float64 `json:"time,omitempty"`
	Watts     *float64 `json:"watts,omitempty"`
	Heartrate *float64 `json:"heartrate,omitempty"`
	Cadence   *float64 `json:"cadence,omitempty"`
}

func (a *activityService) GetActivityStream(_ context.Context, id string) (*ActivityStreamData, error) {
	token, err := a.tokenRepo.Get()
	if err != nil {
		return nil, err
	}
	rawStreams, err := a.storage.GetActivityStream(id)
	if err != nil {
		return nil, err
	}

	if rawStreams == nil {
		rawStreams, err = a.stravaClient.FetchStreams(id, getActivityKeys(), token.AccessToken)
		if err != nil {
			return nil, err
		}
		err = a.storage.SaveActivityStream(id, rawStreams)
		if err != nil {
			return nil, err
		}
	}

	activity, err := a.storage.GetAthleteActivity(id)

	if err != nil {
		return nil, err
	}

	if activity == nil {
		activity, err = a.stravaClient.GetAthleteActivityByID(id, token.AccessToken)
		if err != nil {
			return nil, err
		}
		err = a.storage.SaveAthleteActivity(activity)
		if err != nil {
			return nil, err
		}
	}

	isBike := activity != nil && strings.EqualFold(activity.Type, "ride")

	combined := &ActivityStreamData{
		ActivityID: id,
		Streams:    []StreamDataPoint{},
	}

	if activity != nil {
		combined.ActivityName = activity.Name
		combined.Date = activity.StartDate
	}

	if rawStreams.Time == nil || rawStreams.Watts == nil || rawStreams.Heartrate == nil {
		return combined, nil
	}

	for i := 0; i < len(rawStreams.Time.Data); i++ {
		if rawStreams.Time.Data[i] == nil {
			continue
		}
		if rawStreams.Heartrate.Data[i] == nil {
			continue
		}

		if rawStreams.Watts != nil && rawStreams.Watts.Data[i] == nil && isBike {
			continue
		}
		if rawStreams.Cadence != nil && rawStreams.Cadence.Data[i] == nil && isBike {
			continue
		}
		point := StreamDataPoint{
			Time:      rawStreams.Time.Data[i],
			Heartrate: rawStreams.Heartrate.Data[i],
		}
		if isBike {
			if rawStreams.Cadence != nil {
				point.Cadence = rawStreams.Cadence.Data[i]
			}
			if rawStreams.Watts != nil {
				point.Watts = rawStreams.Watts.Data[i]
			}
		}
		combined.Streams = append(combined.Streams, point)
	}
	return combined, nil
}

func getActivityKeys() []string {
	return []string{"watts", "time", "heartrate", "cadence"}

}
