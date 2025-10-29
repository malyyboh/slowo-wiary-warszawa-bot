package messages

func GetText(key string) string {
	if text, ok := Texts[key]; ok {
		return text
	}
	return Texts["other_answer"]
}

func GetButtonText(buttonSet map[string]string, key string) string {
	if text, ok := buttonSet[key]; ok {
		return text
	}
	return key
}

func HasText(key string) bool {
	_, ok := Texts[key]
	return ok
}
