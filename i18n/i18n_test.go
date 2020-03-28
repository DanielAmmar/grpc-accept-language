package i18n

import (
	"context"
	"encoding/json"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func loadTranslation() {
	enTranslation := `{
	"hello": "Hello world",
	"push": {
        "name1": {
            "title": "example of push title",
            "body": "example of push body"
        }
    }
  }`
	jaTranslation := `{
    "hello": "こんにちは"
  }`

	defaultBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	defaultBundle.MustParseMessageFileBytes([]byte(enTranslation), "en.json")
	defaultBundle.MustParseMessageFileBytes([]byte(jaTranslation), "ja-jp.json")
}

func newMetadataContext(ctx context.Context, val string) context.Context {
	md := metadata.Pairs("accept-language", val)
	return metadata.NewIncomingContext(ctx, md)
}

func TestDefaultLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/test"}
	_, err := UnaryI18nHandler(context.Background(), req, info, func(ctx context.Context, _ interface{}) (interface{}, error) {

		T := MustTfunc(ctx)
		if got, want := T("hello"), "Hello world"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultLanguageOfNestedJson(t *testing.T) {
	loadTranslation()
	req := "request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/test"}
	_, err := UnaryI18nHandler(context.Background(), req, info, func(ctx context.Context, _ interface{}) (interface{}, error) {

		T := MustTfunc(ctx)
		if got, want := T("push.name1.body"), "example of push title"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRespectAcceptLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/test"}
	ctx := newMetadataContext(context.Background(), "ja")
	_, err := UnaryI18nHandler(ctx, req, info, func(ctx context.Context, _ interface{}) (interface{}, error) {

		T := MustTfunc(ctx)
		if got, want := T("hello"), "こんにちは"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFallbackDefaultLanguage(t *testing.T) {
	loadTranslation()
	req := "request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/test"}
	ctx := newMetadataContext(context.Background(), "da")
	_, err := UnaryI18nHandler(ctx, req, info, func(ctx context.Context, _ interface{}) (interface{}, error) {

		T := MustTfunc(ctx)
		if got, want := T("hello"), "Hello world"; got != want {
			t.Errorf("expect T() = %q, but got %q", want, got)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
