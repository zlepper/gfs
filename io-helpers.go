package gfs

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

func isDirectory(p string) (bool, error) {
	stat, err := os.Stat(p)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), err

}

// Gets the statistics for the given directory.
func GetDirectoryStats(fullPath, p string) (*DirectoryStats, error) {
	stats, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	if !stats.IsDir() {
		return nil, errors.New("Cannot get directory stats of something that is not a directory...")
	}

	var dirStats *DirectoryStats
	if p == "/" {
		dirStats = &DirectoryStats{
			Name:                 p,
			Path:                 p,
			LastModificationTime: stats.ModTime(),
		}
	} else {
		dirStats = &DirectoryStats{
			Name:                 stats.Name(),
			Path:                 p,
			LastModificationTime: stats.ModTime(),
		}
	}

	entries, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	dirStats.Entries = make([]DirectoryEntry, len(entries))

	for index, entry := range entries {
		dirEntry := DirectoryEntry{
			Path:                 path.Join(p, entry.Name()),
			Name:                 entry.Name(),
			LastModificationTime: entry.ModTime(),
			IsDirectory:          entry.IsDir(),
		}

		if !entry.IsDir() {
			dirEntry.Size = entry.Size()
		}

		dirStats.Entries[index] = dirEntry
	}

	return dirStats, nil
}
