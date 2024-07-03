package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactSensitiveWithNoCustomFields(t *testing.T) {
	original := struct {
		ApiKey   string
		Username string
	}{"secret", "user"}

	redacted := RedactSensitive(original).(struct {
		ApiKey   string
		Username string
	})

	assert.Equal(t, "REDACTED", redacted.ApiKey)
	assert.Equal(t, "user", redacted.Username)
}

func TestRedactSensitiveWithCustomFields(t *testing.T) {
	original := struct {
		Foo  string
		Spam string
		Bar  string
	}{"some", "weird", "fields"}

	redacted := RedactSensitive(original, "spam").(struct {
		Foo  string
		Spam string
		Bar  string
	})

	assert.Equal(t, "some", redacted.Foo)
	assert.Equal(t, "REDACTED", redacted.Spam)
	assert.Equal(t, "fields", redacted.Bar)
}

func TestRedactSensitiveWithUnsupportedType(t *testing.T) {
	original := "string"
	redacted := RedactSensitive(original)
	assert.Equal(t, original, redacted)
}

func TestRedactSensitiveWithPointerToStruct(t *testing.T) {
	original := &struct {
		Token string
		Info  string
	}{"token", "info"}

	redacted := RedactSensitive(original).(struct {
		Token string
		Info  string
	})

	assert.Equal(t, "REDACTED", redacted.Token)
	assert.Equal(t, "info", redacted.Info)
}

func TestRedactSensitiveWithNestedStruct(t *testing.T) {
	original := struct {
		Credentials struct {
			Password string
		}
		Detail string
	}{struct{ Password string }{"pass"}, "detail"}

	redacted := RedactSensitive(original).(struct {
		Credentials struct {
			Password string
		}
		Detail string
	})

	assert.Equal(t, "REDACTED", redacted.Credentials.Password)
	assert.Equal(t, "detail", redacted.Detail)
}
