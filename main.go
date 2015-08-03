package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/google/go-github/github"
	"github.com/libgit2/git2go"
	"github.com/vpommier/dummy/config"
)

func cloneRepositories(repo *github.Repository, rootDir *string) {
	log.Printf("Cloning repository: %s", *repo.GitURL)

	path := *rootDir + "/" + *repo.Name + ".git"
	opt := &git.CloneOptions{Bare: true}
	_, err := git.Clone(*repo.GitURL, path, opt)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Cloning %s finished", *repo.Name)
}

func getUserRepositories(user *string, client *github.Client, async bool) (allRepos []github.Repository) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	if async == true {
		repos, resp, err := client.Repositories.List(*user, opt)
		if err != nil {
			log.Fatalln(err)
		}
		allRepos = append(allRepos, repos...)

		if resp.NextPage != 0 {
			c := make(chan []github.Repository)
			for i := resp.NextPage; i <= resp.LastPage; i++ {
				opt.ListOptions.Page = i
				go func(opt github.RepositoryListOptions) {
					repos, _, err := client.Repositories.List(*user, &opt)
					if err != nil {
						log.Fatalln(err)
					}
					c <- repos
				}(*opt)
			}
			for i := resp.NextPage; i <= resp.LastPage; i++ {
				allRepos = append(allRepos, <-c...)
			}
		}
	} else {
		for {
			repos, resp, err := client.Repositories.List(*user, opt)
			if err != nil {
				log.Fatalln(err)
			}
			allRepos = append(allRepos, repos...)

			if resp.NextPage == 0 {
				break
			}
			opt.ListOptions.Page = resp.NextPage
		}
	}
	return allRepos
}

func backup(repos []github.Repository, async bool) {
	if async == true {
		var wg sync.WaitGroup
		for i := 0; i < len(repos); i++ {
			wg.Add(1)
			go func(repo *github.Repository) {
				defer wg.Done()
				cloneRepositories(repo, &config.ReposDir)
			}(&repos[i])
		}
		wg.Wait()
	} else {
		for i := 0; i < len(repos); i++ {
			cloneRepositories(&repos[i], &config.ReposDir)
		}
	}
}

func main() {
	user := os.Args[1]
	asyncBackup, err := strconv.ParseBool(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	asyncRequests, err := strconv.ParseBool(os.Args[3])
	if err != nil {
		log.Fatalln(err)
	}

	// Retrieve Github user's repositories.
	t := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     os.Getenv("GITHUB_CLIENTID"),
		ClientSecret: os.Getenv("GITHUB_CLIENTSECRET"),
	}

	client := github.NewClient(t.Client())
	allRepos := getUserRepositories(&user, client, asyncRequests)

	log.Printf("Nb repos: %d", len(allRepos))

	// Clone and compress repositories.
	backup(allRepos, asyncBackup)
}
