package language

import (
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	//langs := []string{"en", "zh-TW", "zh-CN", "ja"}
	langs := []string{"en", "zh-TW"}

	for _, lang := range langs {
		bundle.MustLoadMessageFile(fmt.Sprintf("locales/pages/%v.json", lang))
	}

}

func GetPageMsg(code string, lang string) string {
	localizer := i18n.NewLocalizer(bundle, lang)

	Msg := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: code[0:5],
		},
		TemplateData: map[string]interface{}{
			"Code": code,
		},
	})
	return Msg
}
