package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
)

type I18nMiddleware struct {
	config []string
}

func NewI18nMiddleware(config ...string) *I18nMiddleware {
	return &I18nMiddleware{
		config: config,
	}
}

func (m *I18nMiddleware) Handle(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bundle := bundle(r, m.config...)
		handle(w, WithLanguage(r, bundle))
	}
}

func bundle(r *http.Request, configs ...string) *i18n.Bundle {
	accept := GetAcceptLanguage(r)
	languageTag := language.English
	if accept == "zh-CN" {
		languageTag = language.Chinese
	}
	bundle := i18n.NewBundle(languageTag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, file := range configs {
		bundle.LoadMessageFile(file)
	}
	return bundle
}
