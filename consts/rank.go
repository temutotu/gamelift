package consts

const (
	BIGINNER = 0
	BASIC    = 3000
	EXPERT   = 6000
)

func GetRankName(rank int) string {
	switch rank {
	case BIGINNER:
		return "BIGINNER"
	case BASIC:
		return "BASIC"
	case EXPERT:
		return "EXPERT"
	default:
		return ""
	}
}

func GetRankLabel(rating int) string {
	switch {
	case rating < BASIC:
		return "BIGINNER"
	case rating < EXPERT:
		return "BASIC"
	case rating >= EXPERT:
		return "Expert"
	default:
		return "Basic"
	}
}
