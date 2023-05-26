package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	// Open the repository at the current directory
	r, err := git.PlainOpen("testdata/simplediff")
	if err != nil {
		log.Fatal(err)
	}
	// Get the HEAD reference
	ref, err := r.Head()
	if err != nil {
		log.Fatal(err)
	}

	// Get the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// Get the tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		log.Fatal(err)
	}

	// Get the index (staging area) from the repository
	idx, err := r.Storer.Index()
	if err != nil {
		log.Fatal(err)
	}
	// Iterate over the staged files
	for _, change := range idx.Entries {
		// Get the file path
		filePath := filepath.Join("testdata/simplediff", change.Name)

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Get the object (blob) from the tree
		obj, err := tree.File(change.Name)
		if err != nil {
			log.Fatal(err)
		}
		//cont, _ := obj.Contents()
		br, _ := obj.Blob.Reader()
		bbytes, _ := io.ReadAll(br)
		fbytes, _ := io.ReadAll(file)

		dmp := diffmatchpatch.New()
		fmt.Println("number of \\gls detected before: ", bytes.Count(bbytes, []byte("\\gls")))
		fmt.Println("number of \\gls detected after:  ", bytes.Count(fbytes, []byte("\\gls")))
		diffs := dmp.DiffMain(string(bbytes), string(fbytes), false)
		fmt.Println(dmp.DiffPrettyText(diffs))
	}
}
