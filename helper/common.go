package helper

const (
	CustomInputTextType             = "text"
	UserSettingCardType             = "card"
	UserSettingTimezoneType         = "timezone"
	UserSettingCalendarNotification = "calendar-notification"
	UserSettingMessageNotification  = "message-notification"
	UserSettingActivityNotification = "activity-notification"
	UserSettingEmailNotification    = "email-notification"
)

type UserActionCard struct {
	Name string
}

var DefaultOwnerAdminCards = []UserActionCard{
	{
		Name: "Complete your profile",
	},
	{
		Name: "Invite members",
	},
	{
		Name: "Customize settings",
	},
	{
		Name: "Customize onboarding",
	},
}
