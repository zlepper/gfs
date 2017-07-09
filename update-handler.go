package gfs

import (
	"github.com/zlepper/gfs/internal"
	ghc "github.com/zlepper/github-release-checker"
	"log"
)

var hasUpdate string = ""

func checkForUpdates() {
	release, err := ghc.GetLatestReleaseForPlatform("zlepper", "gfs", internal.FilenameRegex, true)
	if err == nil {
		newer, err := ghc.IsNewer(release, GFSVersion)

		if err != nil {
			log.Println("Error when comparing update versions", err.Error())
		}

		if newer {
			log.Printf("\n\nA newer GFS release is available on github. Download at:\n%s\n\n", release.DownloadUrl)
			hasUpdate = release.DownloadUrl
		} else {
			log.Println("No new version available")
		}
	} else {
		log.Println("Error when checking for updates:", err.Error())
	}
}
