package gfs

import (
	"github.com/zlepper/gfs/internal"
	ghc "github.com/zlepper/github-release-checker"
	"log"
)

var hasUpdate string = ""

func checkForUpdates() {
	release, err := ghc.GetLatestReleaseForPlatform("zlepper", "gfs", internal.FilenameRegex, false)
	if err == nil {
		if newer, err := ghc.IsNewer(release, GFSVersion); err != nil && newer {
			log.Printf("\n\nA newer GFS release is available on github. Download at:\n%s\n\n", release.DownloadUrl)
			hasUpdate = release.DownloadUrl
		}
	}
}
