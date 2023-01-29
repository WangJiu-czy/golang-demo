package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	devMode := flag.Bool("dev", false, "disables adding, commit and pushing, only logs the generated commit message")
	debug := flag.Bool("debug", false, "enable debug default:false")
	conf := getConfig()

	flag.Parse()
	if !CheckForGit(conf) {
		log.Fatalln("[FATAL ERROR] 'git' executable not found,git is required to work properly -exiting.")
	}
	conf.DebugMode = *debug
	if conf.DebugMode {
		DebugLog(conf, "Debug mode enabled")
	}
	if !CheckIfGitRepo(conf) {
		log.Fatalln("[FATAL ERROR] not a git repository -exiting.")
	}
	if *devMode {
		conf.DebugMode = true
		DebugLog(conf, "Dev mode enabled, automatically enabled debug mode, adding, committing and pushing will be disabled.")
		log.Println(generateCommitContent(conf))
		os.Exit(0)
	}
	if conf.PullOnStart {
		log.Println("pulling changes from remote...")
		GitPush(conf)
	}
	log.Println("Watching for changes...")
	for true {
		if GitRepoHashChanges(conf) {
			GitAdd(conf)
			GitCommit(conf)
			GitPush(conf)
			log.Printf("All done, waiting for %d seconds before checking for changes again...", conf.BackupInterval)
		} else {
			log.Println("No changes to commit, waiting for next iteration...")
		}
		time.Sleep(time.Duration(conf.BackupInterval) * time.Second)
	}

}
