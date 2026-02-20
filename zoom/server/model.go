package server

type Notification[T any] struct {
	Event   string `json:"event"`
	EventTs int64  `json:"event_ts"`
	Payload T      `json:"payload"`
}

type webhookHeader struct {
	Event string `json:"event"`
}

type validateTokenPayload struct {
	PlainToken string `json:"plain_token"`
}

type MeetingSettings struct {
	AlternativeHosts string `json:"alternative_hosts"`
	UsePmi           bool   `json:"use_pmi"`
	JbhTime          int    `json:"jbh_time,omitempty"`
	JoinBeforeHost   bool   `json:"join_before_host,omitempty"`
	MeetingInvitees  []struct {
		Email string `json:"email,omitempty"`
	} `json:"meeting_invitees,omitempty"`
}

type MeetingOcurrence struct {
	OccurrenceID string `json:"occurrence_id"`
	StartTime    string `json:"start_time"`
	Duration     int    `json:"duration,omitempty"`
	Status       string `json:"status,omitempty"`
}

type MeetingRecurrence struct {
	EndDateTime    string `json:"end_date_time,omitempty"`
	EndTimes       int    `json:"end_times,omitempty"`
	MonthlyDay     int    `json:"monthly_day,omitempty"`
	MonthlyWeek    int    `json:"monthly_week,omitempty"`
	MonthlyWeekDay int    `json:"monthly_week_day,omitempty"`
	RepeatInterval int    `json:"repeat_interval,omitempty"`
	Type           int    `json:"type,omitempty"`
	WeeklyDays     string `json:"weekly_days,omitempty"`
}

type MeetingObject struct {
	Duration    int                `json:"duration"`
	HostID      string             `json:"host_id"`
	ID          int64              `json:"id"`
	JoinURL     string             `json:"join_url"`
	Settings    MeetingSettings    `json:"settings"`
	Topic       string             `json:"topic"`
	Type        int                `json:"type"`
	UUID        string             `json:"uuid"`
	Occurrences []MeetingOcurrence `json:"occurrence,omitempty"`
	Password    string             `json:"password,omitempty"`
	PMI         string             `json:"pmi,omitempty"`
	Recurrence  MeetingRecurrence  `json:"recurrence"`
	StartTime   string             `json:"start_time,omitempty"`
	Timezone    string             `json:"timezone,omitempty"`
	Issues      []string           `json:"issues,omitempty"`
}

type MeetingEvent struct {
	AccountId  string        `json:"account_id"`
	Object     MeetingObject `json:"object"`
	Operator   string        `json:"operator,omitempty"`
	OperatorId string        `json:"operator_id,omitempty"`
	Operation  string        `json:"operation,omitempty"`
}
