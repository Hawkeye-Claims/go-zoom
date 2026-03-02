package enums

// MeetingType classifies a Zoom meeting by its scheduling behaviour.
type MeetingType int

const (
	// InstantMeeting is a one-time meeting that starts immediately.
	InstantMeeting MeetingType = 1
	// ScheduledMeeting is a one-time meeting scheduled for a future date and time.
	ScheduledMeeting MeetingType = 2
	// RecurringMeetingNoFixedTime is a recurring meeting with no fixed occurrence
	// schedule; the host can start it at any time.
	RecurringMeetingNoFixedTime MeetingType = 3
	// RecurringMeetingFixedTime is a recurring meeting with a fixed, pre-defined
	// occurrence schedule.
	RecurringMeetingFixedTime MeetingType = 8
	// ShareScreenOnly is a meeting in which only screen sharing is enabled.
	ShareScreenOnly MeetingType = 10
)

// MeetingRecurrenceType specifies the cadence of a recurring meeting.
type MeetingRecurrenceType int

const (
	// DailyRecurrence repeats the meeting every day (or every N days).
	DailyRecurrence MeetingRecurrenceType = 1
	// WeeklyRecurrence repeats the meeting every week (or every N weeks).
	WeeklyRecurrence MeetingRecurrenceType = 2
	// MonthlyRecurrence repeats the meeting every month (or every N months).
	MonthlyRecurrence MeetingRecurrenceType = 3
)

// TimeFilterField selects which timestamp field is used when filtering meeting
// summaries by a date range.
type TimeFilterField string

const (
	// SummaryStartTime filters summaries by the time the meeting started.
	SummaryStartTime TimeFilterField = "summary_start_time"
	// SummaryCreatedTime filters summaries by the time the AI summary was created.
	SummaryCreatedTime TimeFilterField = "summary_created"
)
