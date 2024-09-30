package node

import (
	"encoding/json"
	"strconv"
	"time"
)

type Session struct {
	ConsumerCountry string        `json:"consumerCountry"`
	ServiceType     string        `json:"serviceType"`
	Duration        time.Duration `json:"duration"`
	Earning         float64       `json:"earning"`
	Transferred     int64         `json:"transferred"`
	StartedAt       time.Time     `json:"startedAt"`
}

func (s *Session) UnmarshalJSON(data []byte) error {
	type Alias Session
	aux := &struct {
		Earning     string `json:"earning"`
		Transferred string `json:"transferred"`
		Duration    int    `json:"duration"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if aux.Earning != "" {
		s.Earning, err = strconv.ParseFloat(aux.Earning, 64)
		if err != nil {
			return err
		}
	}

	if aux.Transferred != "" {
		s.Transferred, err = strconv.ParseInt(aux.Transferred, 10, 64)
		if err != nil {
			return err
		}
	}

	s.Duration = time.Second * time.Duration(aux.Duration)
	return nil
}
