package domain

type ChannelType string

const (
	EmailChannel ChannelType = "email"
	SMSChannel   ChannelType = "sms"
)

type ChannelMeta struct {
	Name      ChannelType
	Start     int
	End       int
	RegexExpr string
}

var ChannelDefinitions = map[ChannelType]ChannelMeta{
	EmailChannel: {
		Name:      EmailChannel,
		Start:     0,
		End:       40,
		RegexExpr: `^[^\s@]+@[^\s@]+\.[^\s@]+$`,
	},
	SMSChannel: {
		Name:      SMSChannel,
		Start:     41,
		End:       55,
		RegexExpr: `^\d{3}-\d{4}-\d{4}$`,
	},
}
