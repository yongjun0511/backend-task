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
		End:       50,
		RegexExpr: `^[^\s@]+@[^\s@]+\.[^\s@]+$`, // XXXX@XXXX. 형태
	},
	PhoneField: {
		Start:     50,
		End:       71,
		RegexExpr: `^\d{3}-\d{4}-\d{4}$`, // 000-0000-0000 형태
	},
	ScoreUpField: {
		Start:     71,
		End:       72,
		RegexExpr: `^[YN]$`, // Y 또는 N
	},
}
