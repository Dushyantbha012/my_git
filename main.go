package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "usage:":
		fmt.Fprintf(os.Stderr, "will implement it\n")
	case "init":
		//Uncomment this block to pass the first stage!

		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
			os.Exit(1)
		}
		if os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "error, second params needs to be '-p'\n")
			os.Exit(1)
		}
		dirname := os.Args[3][:2]
		filename := os.Args[3][2:]
		objectPath := path.Join(".git", "objects", dirname, filename)
		fileBytes, err := os.ReadFile(objectPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read the file")
			os.Exit(1)
		}
		bytesReader := bytes.NewReader(fileBytes)
		r, err := zlib.NewReader(bytesReader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating zlib new reader")
			os.Exit(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read the file")
			os.Exit(1)
		}
		defer r.Close()
		var decompressedData bytes.Buffer
		_, err = decompressedData.ReadFrom(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error decompressing data")
			os.Exit(1)
		}
		resultString := decompressedData.String()
		splits := strings.SplitN(resultString, "\x00", 2)
		_, content := splits[0], splits[1]
		fmt.Print(content)
	case "hash-object":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
			os.Exit(1)
		} else if len(os.Args) == 3 {
			filePath := os.Args[2]
			file, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file \n hello")
				print(err)
				os.Exit(1)
			}
			content := string(file)
			stats, err := os.Stat(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading stats ")
				os.Exit(1)
			}
			contentAndHeader := fmt.Sprintf("blob %d\x00%s", stats.Size(), content)
			sha := sha1.Sum([]byte(contentAndHeader))
			hash := fmt.Sprintf("%x", sha)
			println(hash)
		} else if len(os.Args) == 4 {
			cmd := os.Args[2]
			if cmd != "-w" {
				fmt.Fprintf(os.Stderr, "error, second params needs to be '-w'\n")
				os.Exit(1)
			}
			filePath := os.Args[3]
			println("Filepath: %s\n", filePath)
			file, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file\n")
				fmt.Println(err)
				os.Exit(1)
			}
			content := string(file)
			stats, err := os.Stat(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading stats")
				os.Exit(1)
			}
			contentAndHeader := fmt.Sprintf("blob %d\x00%s", stats.Size(), content)
			sha := sha1.Sum([]byte(contentAndHeader))
			hash := fmt.Sprintf("%x", sha)
			blobName := []rune(hash)
			blobPath := ".git/objects/"
			for i, v := range blobName {
				blobPath += string(v)
				if i == 1 {
					blobPath += "/"
				}
			}
			var buffer bytes.Buffer
			z := zlib.NewWriter(&buffer)
			z.Write([]byte(contentAndHeader))
			z.Close()
			os.MkdirAll(filepath.Dir(blobPath), os.ModePerm)
			f, _ := os.Create(blobPath)
			defer f.Close()
			f.Write(buffer.Bytes())
			fmt.Print(hash)
		} else {
			fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
