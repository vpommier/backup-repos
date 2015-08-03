package config

import (
	"os"
)

var (
	ReposDir    string
	ArchivesDir string
)

var defaultReposDir = "/var/backup-repos/repos"
var defaultArchivesDir = "/var/backup-repos/archives"

func init() {
	ReposDir = os.Getenv("BACKUP_REPOS_REPOS_DIR")
	ArchivesDir = os.Getenv("BACKUP_REPOS_ARCHIVES_DIR")
	if ReposDir == "" {
		ReposDir = defaultReposDir
	}
	if ArchivesDir == "" {
		ArchivesDir = defaultArchivesDir
	}
}
