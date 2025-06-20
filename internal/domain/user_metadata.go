package domain

type FieldType string

const (
	EmailField   FieldType = "email"
	PhoneField   FieldType = "phone"
	ScoreUpField FieldType = "score_up"
)

type FieldMeta struct {
	Start     int
	End       int
	RegexExpr string
}

var UserFieldDefinitions = map[FieldType]FieldMeta{
	EmailField: {
		Start:     0,
		End:       40,
		RegexExpr: `^[^\s@]+@[^\s@]+\.[^\s@]+$`,
	},
	PhoneField: {
		Start:     41,
		End:       59,
		RegexExpr: `^\d{3}-\d{4}-\d{4}$`,
	},
	ScoreUpField: {
		Start:     60,
		End:       61,
		RegexExpr: `^[YN]$`,
	},
}
