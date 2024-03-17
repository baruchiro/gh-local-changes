package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"

	"github.com/charmbracelet/log"
)

var skipFolders = []string{
	"node_modules",
}
var gitReposChan = make(chan string, 10)

type results struct {
	repo     string
	branches map[string]string
	changes  int
}

var allResults = []results{}

func walkDirectory(dir string) error {
	defer close(gitReposChan)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			for _, skip := range skipFolders {
				if info.Name() == skip {
					return filepath.SkipDir
				}
			}
			if info.Name() == ".git" {
				repoFolder := filepath.Dir(path)
				log.Debug("Found git repository:", repoFolder)
				gitReposChan <- repoFolder
				return filepath.SkipDir
			}
		}
		return nil
	})
}

func handleGitRepos(wg *sync.WaitGroup) {
	defer wg.Done()
	for repoFolder := range gitReposChan {
		log.Debug("Processing git repository:", repoFolder)
		repo := &GitRepo{folder: repoFolder}
		results := results{
			repo:     repoFolder,
			branches: map[string]string{},
			changes:  0,
		}

		branches, err := repo.getUnpushedBranches()
		if err != nil {
			log.Error("Error checking git status:", err)
			continue
		} else {
			results.branches = branches
		}

		changes, err := repo.getUnpushedChanges()
		if err != nil {
			log.Error("Error checking git status:", err)
			continue
		} else {
			results.changes = changes
		}

		allResults = append(allResults, results)
	}
}

func printLookingForContributors() {
	log.Info("--------------------------------------------")
	log.Info("This project was developed in a short time, before I returned my laptop to the IT department. I'm looking for contributors to help me improve it.")
	log.Info("Any help is welcome, be it a suggestion, a bug report, a pull request...")
	log.Info("https://github.com/baruchiro/gh-local-changes")
	log.Info("--------------------------------------------")
}

func main() {
	defer printLookingForContributors()
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
		if dir == "-h" || dir == "--help" {
			log.Info("Usage: go-gh [dir]")
			os.Exit(0)
		}
		// Fail if the directory does not exist
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatal("Directory does not exist:", dir)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go handleGitRepos(wg)

	if err := walkDirectory(dir); err != nil {
		log.Fatal("Error walking directory:", err)
		os.Exit(1)
	}
	wg.Wait()

	// Print results in a table
	log.Info("Results:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Repo\tBranches\tChanges")
	for _, r := range allResults {
		if len(r.branches) > 0 || r.changes > 0 {
			fmt.Fprintf(w, "%s\t%v\t%d\n", r.repo, len(r.branches), r.changes)
		}
	}
	w.Flush()
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
