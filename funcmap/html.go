package funcmap

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/url"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func addHTMLFuncs(target map[string]any) {

	target["addQueryParams"] = func(extraParams string, url string) string {
		if strings.Contains(url, "?") {
			return url + "&" + extraParams
		}
		return url + "?" + extraParams
	}

	target["attr"] = func(value string) template.HTMLAttr {
		return template.HTMLAttr(value)
	}

	target["css"] = func(value string) template.CSS {
		return template.CSS(value)
	}

	target["highlight"] = func(text string, search string) template.HTML {
		if search == "" {
			return template.HTML(text)
		}
		result := strings.ReplaceAll(text, search, `<b class="highlight">`+search+"</b>")
		return template.HTML(result)
	}

	target["domainOnly"] = func(value string) string {
		result := strings.TrimPrefix(value, "http://")
		result = strings.TrimPrefix(result, "https://")
		result, _, _ = strings.Cut(result, "/")
		return result
	}

	target["hasImage"] = func(value string) bool {
		if strings.Contains(value, "<img") {
			return true
		}

		if strings.Contains(value, "<picture") {
			return true
		}

		return false
	}

	target["html"] = func(value string) template.HTML {
		return template.HTML(value)
	}

	target["htmlMinimal"] = func(value string) template.HTML {
		return template.HTML(html.Minimal(value))
	}

	target["js"] = func(value string) string {
		return template.JSEscapeString(value)
	}

	target["json"] = func(value any) string {
		result, _ := json.Marshal(value)
		return string(result)
	}

	target["jsonIndent"] = func(value any) string {
		result, _ := json.MarshalIndent(value, "", "    ")
		return string(result)
	}

	target["markdown"] = func(value any) template.HTML {

		valueBytes := convert.Bytes(value)

		// https://github.com/yuin/goldmark#built-in-extensions
		var buffer bytes.Buffer

		md := goldmark.New(
			goldmark.WithExtensions(
				extension.Table,
				extension.Linkify,
				extension.Typographer,
				extension.DefinitionList,
			),
			goldmark.WithRendererOptions(),
		)

		if err := md.Convert([]byte(valueBytes), &buffer); err != nil {
			derp.Report(derp.Wrap(err, "tools.templates.functions.markdown", "Error converting Markdown to HTML"))
		}

		return template.HTML(buffer.String())
	}

	target["queryEscape"] = func(value string) string {
		return url.QueryEscape(value)
	}

	target["removeLinks"] = func(value string) template.HTML {
		result := strings.ReplaceAll(value, "<a ", "<span ")
		result = strings.ReplaceAll(result, "</a", "</span")
		return template.HTML(result)
	}

	target["summary"] = html.Summary

	target["stripProtocol"] = func(value string) string {
		if after, found := strings.CutPrefix(value, "http://"); found {
			return after
		}
		if after, found := strings.CutPrefix(value, "https://"); found {
			return after
		}
		return value
	}

	target["textOnly"] = html.RemoveTags

	target["text"] = func(value string) template.HTML {
		return template.HTML(html.FromText(value))
	}
}
