package model

import "time"

type TokenRequest struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (t *TokenResponse) IsExpired() bool {
	now := time.Now().Unix()
	return t.ExpiresAt <= now
}

type RedirectTokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	Athlete      struct {
		ID            int     `json:"id,omitempty"`
		Username      string  `json:"username,omitempty"`
		ResourceState int     `json:"resource_state,omitempty"`
		Firstname     string  `json:"firstname,omitempty"`
		Lastname      string  `json:"lastname,omitempty"`
		Bio           string  `json:"bio,omitempty"`
		City          string  `json:"city,omitempty"`
		State         string  `json:"state,omitempty"`
		Country       string  `json:"country,omitempty"`
		Sex           string  `json:"sex,omitempty"`
		Premium       bool    `json:"premium,omitempty"`
		Summit        bool    `json:"summit,omitempty"`
		CreatedAt     string  `json:"created_at,omitempty"`
		UpdatedAt     string  `json:"updated_at,omitempty"`
		BadgeTypeID   int     `json:"badge_type_id,omitempty"`
		Weight        float64 `json:"weight,omitempty"`
		ProfileMedium string  `json:"profile_medium,omitempty"`
		Profile       string  `json:"profile,omitempty"`
	} `json:"athlete,omitempty"`
}

func (t *RedirectTokenResponse) IsExpired() bool {
	now := time.Now().Unix()
	return t.ExpiresAt <= now
}

type Athlete struct {
	ID            int64 `json:"id,omitempty"`
	ResourceState int   `json:"resource_state,omitempty"`
}

type Map struct {
	ID              string `json:"id,omitempty"`
	SummaryPolyline string `json:"summary_polyline,omitempty"`
	ResourceState   int    `json:"resource_state,omitempty"`
}

type AthleteActivity struct {
	ResourceState              int       `json:"resource_state,omitempty"`
	Athlete                    Athlete   `json:"athlete,omitempty"`
	Name                       string    `json:"name,omitempty"`
	Distance                   float64   `json:"distance,omitempty"`
	MovingTime                 int       `json:"moving_time,omitempty"`
	ElapsedTime                int       `json:"elapsed_time,omitempty"`
	TotalElevationGain         float64   `json:"total_elevation_gain,omitempty"`
	Type                       string    `json:"type,omitempty"`
	SportType                  string    `json:"sport_type,omitempty"`
	WorkoutType                *int      `json:"workout_type,omitempty"`
	ID                         int64     `json:"id,omitempty"`
	StartDate                  string    `json:"start_date,omitempty"`
	StartDateLocal             string    `json:"start_date_local,omitempty"`
	Timezone                   string    `json:"timezone,omitempty"`
	UTCOffset                  float64   `json:"utc_offset,omitempty"`
	LocationCity               *string   `json:"location_city,omitempty"`
	LocationState              *string   `json:"location_state,omitempty"`
	LocationCountry            *string   `json:"location_country,omitempty"`
	AchievementCount           int       `json:"achievement_count,omitempty"`
	KudosCount                 int       `json:"kudos_count,omitempty"`
	CommentCount               int       `json:"comment_count,omitempty"`
	AthleteCount               int       `json:"athlete_count,omitempty"`
	PhotoCount                 int       `json:"photo_count,omitempty"`
	Map                        Map       `json:"map,omitempty"`
	Trainer                    bool      `json:"trainer,omitempty"`
	Commute                    bool      `json:"commute,omitempty"`
	Manual                     bool      `json:"manual,omitempty"`
	Private                    bool      `json:"private,omitempty"`
	Visibility                 string    `json:"visibility,omitempty"`
	Flagged                    bool      `json:"flagged,omitempty"`
	GearID                     *string   `json:"gear_id,omitempty"`
	StartLatLng                []float64 `json:"start_latlng,omitempty"`
	EndLatLng                  []float64 `json:"end_latlng,omitempty"`
	AverageSpeed               float64   `json:"average_speed,omitempty"`
	MaxSpeed                   float64   `json:"max_speed,omitempty"`
	AverageCadence             *float64  `json:"average_cadence,omitempty"`
	AverageTemp                *int      `json:"average_temp,omitempty"`
	AverageWatts               *float64  `json:"average_watts,omitempty"`
	MaxWatts                   *int      `json:"max_watts,omitempty"`
	WeightedAverageWatts       *int      `json:"weighted_average_watts,omitempty"`
	DeviceWatts                *bool     `json:"device_watts,omitempty"`
	Kilojoules                 *float64  `json:"kilojoules,omitempty"`
	HasHeartrate               bool      `json:"has_heartrate,omitempty"`
	AverageHeartrate           *float64  `json:"average_heartrate,omitempty"`
	MaxHeartrate               *float64  `json:"max_heartrate,omitempty"`
	HeartrateOptOut            bool      `json:"heartrate_opt_out,omitempty"`
	DisplayHideHeartrateOption bool      `json:"display_hide_heartrate_option,omitempty"`
	ElevHigh                   *float64  `json:"elev_high,omitempty"`
	ElevLow                    *float64  `json:"elev_low,omitempty"`
	UploadID                   int64     `json:"upload_id,omitempty"`
	UploadIDStr                string    `json:"upload_id_str,omitempty"`
	ExternalID                 string    `json:"external_id,omitempty"`
	FromAcceptedTag            bool      `json:"from_accepted_tag,omitempty"`
	PRCount                    int       `json:"pr_count,omitempty"`
	TotalPhotoCount            int       `json:"total_photo_count,omitempty"`
	HasKudoed                  bool      `json:"has_kudoed,omitempty"`
	SufferScore                *float64  `json:"suffer_score,omitempty"`
}

type StreamData struct {
	Data         []*float64 `json:"data,omitempty"`
	SeriesType   string     `json:"series_type,omitempty"`
	OriginalSize int32      `json:"original_size,omitempty"`
	Resolution   string     `json:"resolution,omitempty"`
}

type ActivityStreams struct {
	Watts     *StreamData `json:"watts,omitempty"`
	Time      *StreamData `json:"time,omitempty"`
	Heartrate *StreamData `json:"heartrate,omitempty"`
	Cadence   *StreamData `json:"cadence,omitempty"`
}
