package enums

type MeetingType int

const (
	InstantMeeting              MeetingType = 1
	ScheduledMeeting            MeetingType = 2
	RecurringMeetingNoFixedTime MeetingType = 3
	RecurringMeetingFixedTime   MeetingType = 8
)

type MeetingRecurrenceType int

const (
	DailyRecurrence   MeetingRecurrenceType = 1
	WeeklyRecurrence  MeetingRecurrenceType = 2
	MonthlyRecurrence MeetingRecurrenceType = 3
)
