#!/bin/bash

usage(){
        echo "Usage: backup-repos.sh <provider username>"
        echo "Example: backup-repos.sh vpommier"
}

log(){
	echo "$(date +'%D %T') INFO: $*" 1>&2
}

error(){
	echo "$(date +'%D %T') ERROR: $*" 1>&2
}

getAuth(){
	if [[ -z "${CLIENT_ID}" || -z "${CLIENT_SECRET}" ]]
	then
		echo ""
	else
		echo "client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&"
	fi
}

getRepos(){
	local i=1
	while true
	do
		curl --fail \
			--silent \
			--show-error \
			--location \
			--dump-header "$headerFile" \
			--OUTPUT "${reposFile}${i}" \
			--url "https://api.github.com/users/${1}/repos?$(getAuth)page=${i}&per_page=100" || return 1
		let i++
		grep -q '^Link: .*rel="next"' "$headerFile" || break
	done
}

backup(){
	for f in "$reposFile"*
	do
		if [ -f "$f" ]
		then
			for r in $(jq .[].git_url "$f")
			do
				local url
				local repoFullName
				local repoName
				url=$(expr "$r" : '"\(.*\)"')
				repoFullName=$(basename "$url")
				repoName=$(expr "$repoFullName" : '\(.*\)\.git')
				log "Cloning repository: $url"
				git clone --quiet --bare "$url" "${reposDir}/$repoFullName"
				log "Cloning $repoFullName finished"
				log "Archiving repository: $repoFullName"
				tar -czf "${BACKUP_REPOS_ARCHIVES_DIR}/${repoName}_$(date +%s).tar" -C "$reposDir" "$repoFullName"
				log "Archiving $repoFullName finished"
			done
		fi
	done
}

cleanTempFiles(){
	rm -rf "$headerFile" "$reposFile"* "$reposDir"
}

if [ -z "$1" ]
then
        usage
        exit 1
fi

[[ $1 =~ ^[a-zA-Z0-9]+$ ]] || {
        error 'Invalid charaters, must be alphanumeric.'
        exit 1
}

BACKUP_REPOS_ARCHIVES_DIR=${BACKUP_REPOS_ARCHIVES_DIR:-/var/backup-repos/archives}

reposDir=/tmp/repos
headerFile=/tmp/header
reposFile=/tmp/repos

mkdir -p "$BACKUP_REPOS_ARCHIVES_DIR"

# Clean temp files
cleanTempFiles

# Do backups
getRepos "$1" || exit 1
backup

# Clean temp files
cleanTempFiles

# Clean old archives
find "$BACKUP_REPOS_ARCHIVES_DIR" -name \*.tar -mtime +7 -exec rm -f {} \;
