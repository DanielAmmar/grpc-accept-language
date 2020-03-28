package i18n

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"google.golang.org/grpc"

	acceptlang "github.com/kazegusuri/grpc-accept-language"
)

var defaultLanguage = "en"

func SetDefaultLanguage(lang string) {
	defaultLanguage = lang
}

var _ grpc.UnaryServerInterceptor = UnaryI18nHandler

type tfuncKey struct{}

var defaultBundle = i18n.NewBundle(language.English)

func UnaryI18nHandler(origctx context.Context, origreq interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return acceptlang.UnaryAcceptLanguageHandler(origctx, origreq, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		acceptLangs := acceptlang.FromContext(ctx)
		langs := acceptLangs.Languages()
		tfunc := MustTfuncLocales(langs...)
		ctx = context.WithValue(ctx, tfuncKey{}, tfunc)
		return handler(ctx, req)
	})
}

func MustTfunc(ctx context.Context) func(string) string {
	tfunc, ok := ctx.Value(tfuncKey{}).(TranslateFunc)
	if !ok {
		panic("could not find TranslateFunc from context")
	}
	return tfunc
}

type TranslateFunc func(key string) string

func MustTfuncLocales(localeIDs ...string) TranslateFunc {
	localizer := i18n.NewLocalizer(defaultBundle, localeIDs...)
	return func(key string) string {
		str, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: key,
		})

		if err != nil {
			return ""
		}

		return str

	}
}
