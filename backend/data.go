package backend

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

var (
	regexIsHttp = regexp.MustCompile(`^https?://`)
)

type Data struct {
	*Config
	Description string
	Title       string
	URL         string
}

func (d *Data) HandleGenericLater(name string, title, desc string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := r.Context().Value(configKey).(*Config)
		d.Config = cfg
		d.URL = "/" + strings.TrimPrefix(r.URL.Path, "/")
		d.Title = fmt.Sprintf("%s - %s", cfg.Name, title)
		d.Description = desc
		t, err := template.New("").Funcs(template.FuncMap{
			"static": getStatic,
			"fullStatic": func(path string) string {
				if regexIsHttp.MatchString(path) {
					return path
				}
				return fmt.Sprintf("https://%s/static/%s", cfg.Domain, path)
			},
			"asset": func(path string) *assetData {
				return getAsset(r.Context(), path)
			},
			"next":   func(i int) int { return i + 1 },
			"before": func(i int) int { return i - 1 },
			"modulo": func(i int, j int) int { return i % j },
		}).ParseFS(templates, fmt.Sprintf("templates/%s.html", name), "templates/base.html")
		if err != nil {
			panic(err)
		}
		exec := "base.html"
		err = t.ExecuteTemplate(w, exec, d)
		if err != nil {
			panic(err)
		}
	}
}

func getStatic(path string) string {
	if regexIsHttp.MatchString(path) {
		return path
	}
	return fmt.Sprintf("/static/%s", path)
}

type assetData struct {
	Src      string
	Checksum string
}

var assets = map[string]*assetData{}

func getAsset(ctx context.Context, path string) *assetData {
	asset, ok := assets[path]
	if ok && !ctx.Value(debugKey).(bool) {
		return asset
	}
	asset = &assetData{}
	var b []byte
	if regexIsHttp.MatchString(path) {
		asset.Src = path
		resp, err := http.Get(path)
		if err != nil {
			slog.Warn("get remote asset", "error", err)
			return asset
		}
		defer resp.Body.Close()
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			slog.Warn("read remote asset", "error", err)
			return asset
		}
	} else {
		asset.Src = fmt.Sprintf("/assets/%s", path)
		aFS := ctx.Value(assetsFSKey).(fs.FS)
		var err error
		b, err = fs.ReadFile(aFS, path)
		if err != nil {
			slog.Warn("read asset", "error", err)
			return asset
		}
	}
	sum := sha256.Sum256(b)
	checksum := base64.StdEncoding.EncodeToString(sum[:])
	asset.Checksum = fmt.Sprintf("sha256-%s", checksum)
	assets[path] = asset
	return asset
}
