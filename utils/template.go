package utils

import (
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

// get HTML i18n string
func i18nHTML(lang, format string, args ...interface{}) template.HTML {
	return template.HTML(i18n.Tr(lang, format, args...))
}

func boolicon(b bool) (s template.HTML) {
	if b {
		s = `<i style="color:green;" class="icon-check""></i>`
	} else {
		s = `<i class="icon-check-empty""></i>`
	}
	return
}

func date(t time.Time) string {
	return beego.Date(t, DateFormat)
}

func datetime(t time.Time) string {
	return beego.Date(t, DateTimeFormat)
}

func loadtimes(t time.Time) int {
	return int(time.Now().Sub(t).Nanoseconds() / 1e6)
}

func sum(base interface{}, value interface{}, params ...interface{}) (s string) {
	switch v := base.(type) {
	case string:
		s = v + ToStr(value)
		for _, p := range params {
			s += ToStr(p)
		}
	}
	return s
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func timesince(lang string, t time.Time) string {
	now := time.Now()
	seonds := int(now.Sub(t).Seconds())
	if seonds < 60 {
		return i18n.Tr(lang, "seconds_ago", seonds)
	} else if seonds < 60*60 {
		return i18n.Tr(lang, "minutes_ago", seonds/60)
	} else if seonds < 60*60*24 {
		return i18n.Tr(lang, "hours_ago", seonds/(60*60))
	} else {
		return i18n.Tr(lang, "days_ago", seonds/(60*60*24))
	}
}

// create an login url with specify redirect to param
func loginto(uris ...string) template.HTMLAttr {
	var uri string
	if len(uris) > 0 {
		uri = uris[0]
	}
	to := fmt.Sprintf("%slogin", AppUrl)
	if len(uri) > 0 {
		to += "?to=" + url.QueryEscape(uri)
	}
	return template.HTMLAttr(to)
}

func generateInitialAjaxUrlVariable(initialUrl string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<script>
			var initialAjaxUrl = '%s';
		</script>`,
		template.HTMLEscapeString(initialUrl)))
}

func init() {
	// Register template functions.
	beego.AddFuncMap("i18n", i18nHTML)
	beego.AddFuncMap("boolicon", boolicon)
	beego.AddFuncMap("date", date)
	beego.AddFuncMap("datetime", datetime)
	beego.AddFuncMap("dict", dict)
	beego.AddFuncMap("timesince", timesince)
	beego.AddFuncMap("loadtimes", loadtimes)
	beego.AddFuncMap("sum", sum)
	beego.AddFuncMap("loginto", loginto)
	beego.AddFuncMap("initajaxurl", generateInitialAjaxUrlVariable)
}
