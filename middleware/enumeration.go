package middleware

func ProductStatusEnum(statusEnum int) string {
	if statusEnum == 1 {
		return "available"
	} else if statusEnum == 2 {
		return "sold"
	} else {
		return "unavailable"
	}
}
