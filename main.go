package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/nritholtz/stdemuxerhook"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var opts struct {
	ModulePath string `short:"p" long:"module_path" default:"./.terrafile/vendor" description:"File path to install generated terraform modules"`

	TerrafilePath string `short:"f" long:"terrafile_file" default:"./.terrafile/Terrafile" description:"File path to the Terrafile file"`
	Debug         bool   `short:"d" long:"debug"`
}

// To be set by goreleaser on build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var tempDir string

func init() {
	// Needed to redirect logrus to proper stream STDOUT vs STDERR
	log.AddHook(stdemuxerhook.New(log.StandardLogger()))
	var err error
	tempDir, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatalln(err)
	}
}

func gitClone(repositoryPath string) string {
	pathParts := strings.Split(repositoryPath, ":")
	repositoryName := pathParts[1]

	repoPath := fmt.Sprintf("%s/%s", tempDir, repositoryName)
	cmd := exec.Command("git", "clone", repositoryPath, repoPath)
	if opts.Debug {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return repoPath
}

func gitCheckoutRef(repositoryPath string, ref string, destinationDir string) {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = repositoryPath
	if opts.Debug {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	destWithSlash := fmt.Sprintf("%s/", destinationDir)
	cmd = exec.Command("git", "checkout-index", "--prefix", destWithSlash, "-a")
	cmd.Dir = repositoryPath
	if opts.Debug {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	//fmt.Printf("Terrafile: version %v, commit %v, built at %v \n", version, commit, date)
	_, err := flags.Parse(&opts)

	// Invalid choice
	if err != nil {
		panic("invalid arguments")
	}

	// Read File
	yamlFile, err := ioutil.ReadFile(opts.TerrafilePath)
	if err != nil {
		panic(err)
	}

	// Parse File
	var config map[string][]string
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}

	// Clone modules
	os.RemoveAll(opts.ModulePath)
	os.MkdirAll(opts.ModulePath, os.ModePerm)
	for source, refs := range config {
		fmt.Printf("[*] Cloning   %s\n", source)
		repo := gitClone(source)
		for _, ref := range refs {
			pathParts := strings.Split(source, ":")
			repositoryName := pathParts[1]
			fmt.Printf("[*] Vendoring ref %s\n", ref)
			targetPath, err := filepath.Abs(fmt.Sprintf("%s/%s/refs/%s", opts.ModulePath, repositoryName, ref))
			if err != nil {
				panic(err)
			}
			gitCheckoutRef(repo, ref, targetPath)
		}
	}
}
