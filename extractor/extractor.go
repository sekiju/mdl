package extractor

import (
	"fmt"
	"github.com/sekiju/mdl/extractor/cmoa"
	"github.com/sekiju/mdl/extractor/comic_walker"
	"github.com/sekiju/mdl/extractor/corocoro"
	"github.com/sekiju/mdl/extractor/giga_viewer"
	"github.com/sekiju/mdl/extractor/storia_takeshobo"
	"github.com/sekiju/mdl/internal/config"
	"github.com/sekiju/mdl/sdk/manga"
)

type Factory func(cookieString *string) manga.Extractor

var registry = map[string]Factory{
	"comic-walker.com": func(session *string) manga.Extractor {
		return comic_walker.New()
	},
	"shonenjumpplus.com":        gigaViewer("shonenjumpplus.com"),
	"comic-zenon.com":           gigaViewer("comic-zenon.com"),
	"pocket.shonenmagazine.com": gigaViewer("pocket.shonenmagazine.com"),
	"comic-gardo.com":           gigaViewer("comic-gardo.com"),
	"magcomi.com":               gigaViewer("magcomi.com"),
	"tonarinoyj.jp":             gigaViewerNoCookies("tonarinoyj.jp"),
	"comic-ogyaaa.com":          gigaViewerNoCookies("comic-ogyaaa.com"),
	"comic-action.com":          gigaViewer("comic-action.com"),
	"comic-days.com":            gigaViewer("comic-days.com"),
	"comic-growl.com":           gigaViewerNoCookies("comic-growl.com"),
	"comic-earthstar.com":       gigaViewerNoCookies("comic-earthstar.com"),
	"comicborder.com":           gigaViewerNoCookies("comicborder.com"),
	"comic-trail.com":           gigaViewerNoCookies("comic-trail.com"),
	"kuragebunch.com":           gigaViewer("kuragebunch.com"),
	"viewer.heros-web.com":      gigaViewer("viewer.heros-web.com"),
	"www.sunday-webry.com":      gigaViewerNoCookies("www.sunday-webry.com"),
	"www.cmoa.jp": func(cookieString *string) manga.Extractor {
		if cookieString != nil {
			return cmoa.New(*cookieString)
		}
		return nil
	},
	"www.corocoro.jp": func(cookieString *string) manga.Extractor {
		if cookieString != nil {
			return corocoro.NewAuthorized(*cookieString)
		}
		return corocoro.New()
	},
	"storia.takeshobo.co.jp": func(cookieString *string) manga.Extractor {
		return storia_takeshobo.New()
	},
}

func gigaViewer(hostname string) Factory {
	return func(cookieString *string) manga.Extractor {
		if cookieString != nil {
			return giga_viewer.NewAuthorized(hostname, *cookieString)
		}
		return giga_viewer.New(hostname)
	}
}

func gigaViewerNoCookies(hostname string) Factory {
	return func(cookieString *string) manga.Extractor {
		return giga_viewer.New(hostname)
	}
}

func NewExtractor(cfg *config.Config, hostname string) (manga.Extractor, error) {
	factory, exists := registry[hostname]
	if !exists {
		return nil, fmt.Errorf("unsupported website: %s", hostname)
	}

	session := getSession(cfg, hostname)
	return factory(session), nil
}

func getSession(cfg *config.Config, hostname string) *string {
	if cfg.PrimaryCookie != nil {
		return cfg.PrimaryCookie
	}

	if site, exists := cfg.Sites[hostname]; exists && site.CookieString != nil {
		return site.CookieString
	}
	return nil
}