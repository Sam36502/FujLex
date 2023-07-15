package view

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	NTC_PARAM        = "ntc"
	NTC_TYPE_INFO    = "info"
	NTC_TYPE_ERROR   = "error"
	NTC_TYPE_WARNING = "warn"
	NTC_TYPE_SUCCESS = "succ"
)

type Notice struct {
	Type  string
	Title string
	Text  string
}

func NoticeParam(kind, title, text string) string {
	enc := encodeNotice(Notice{
		Type:  kind,
		Title: title,
		Text:  text,
	})
	return fmt.Sprintf("%s=%s", NTC_PARAM, enc)
}

func ErrorNotice(title, text string) string   { return NoticeParam(NTC_TYPE_ERROR, title, text) }
func InfoNotice(title, text string) string    { return NoticeParam(NTC_TYPE_INFO, title, text) }
func SuccessNotice(title, text string) string { return NoticeParam(NTC_TYPE_SUCCESS, title, text) }

func GetNotice(c echo.Context) *Notice {
	enc := c.QueryParam(NTC_PARAM)
	if enc == "" {
		return nil
	}
	return decodeNotice(enc)
}

func encodeNotice(n Notice) string {
	escapedTitle := strings.ReplaceAll(n.Title, ";", "&s&")
	escapedText := strings.ReplaceAll(n.Text, ";", "&s&")
	base := fmt.Sprintf("%s;%s;%s", n.Type, escapedTitle, escapedText)
	return base64.StdEncoding.EncodeToString([]byte(base))
}

func decodeNotice(s string) *Notice {
	base, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil
	}
	fields := strings.Split(string(base), ";")
	unescapedTitle := strings.ReplaceAll(fields[1], "&s&", ";")
	unescapedText := strings.ReplaceAll(fields[2], "&s&", ";")
	return &Notice{
		Type:  fields[0],
		Title: unescapedTitle,
		Text:  unescapedText,
	}
}
