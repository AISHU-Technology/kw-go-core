package i18n

import (
	"context"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

const (
	ctxLanguage   string = "I18nLanguage"
	languageOther        = "other"
)

func GetAcceptLanguage(r *http.Request) string {
	s := r.Header.Get("Accept-Language")
	//lang := r.FormValue("lang")
	return s
}

func WithLanguage(r *http.Request, bundle *i18n.Bundle) *http.Request {
	accept := GetAcceptLanguage(r)
	localizer := i18n.NewLocalizer(bundle, accept)
	return r.WithContext(setLocalizer(r.Context(), localizer))
}

func getLocalizer(ctx context.Context) (*i18n.Localizer, bool) {
	v := ctx.Value(ctxLanguage)
	if l, b := v.(*i18n.Localizer); b {
		return l, true
	}
	return nil, false
}

func setLocalizer(ctx context.Context, l *i18n.Localizer) context.Context {
	return context.WithValue(ctx, ctxLanguage, l)
}

func FormatText(ctx context.Context, key string) string {
	return FormatTextWithData(ctx, key, nil)
}
func FormatMessage(ctx context.Context, message *i18n.Message, args map[string]interface{}) string {
	if localizer, ok := getLocalizer(ctx); ok {
		return localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: message,
			TemplateData:   args,
		})
	}
	return formatInternalMessage(message, args)
}

func FormatTextWithData(ctx context.Context, key string, args map[string]interface{}) string {
	return FormatMessage(ctx, &i18n.Message{
		ID: key,
	}, args)
}

func formatInternalMessage(message *i18n.Message, args map[string]interface{}) string {
	if args == nil {
		return message.Other
	}
	tpl := i18n.NewMessageTemplate(message)
	msg, err := tpl.Execute(languageOther, args, nil)
	if err != nil {
		panic(err)
	}
	return msg
}
