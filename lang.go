package porgs

func IsLangSupported(langID string) bool {
	for _, lang := range SiteConfig.LangSupported {
		if lang == langID {
			return true
		}
	}
	return false
}
