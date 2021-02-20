package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var usageStr = `
Usage: fasta-merge [options] FILE1 [FILE2] ...
Options:
     -o --output <filepath>     file path to output result, default: ./merge.fasta
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

//ACTGBaseSequence ACTG base sequence struct
type ACTGBaseSequence struct {
	Name    string
	Content string
}

//Append Append content to target sequence
func (s *ACTGBaseSequence) Append(newContent string) {
	s.Content += newContent
}

//NewACTG create new ACTG
func NewACTG(name, content string) *ACTGBaseSequence {
	return &ACTGBaseSequence{
		Name:    name,
		Content: content,
	}
}

func main() {
	fs := flag.NewFlagSet("f", flag.ExitOnError)
	fs.Usage = usage

	var outputFile string

	fs.StringVar(&outputFile, "o", "./merge.fasta", "")
	fs.StringVar(&outputFile, "output", "./merge.fasta", "")
	fs.Parse(os.Args[1:])
	fileNames := fs.Args()
	if len(fileNames) <= 0 {
		usage()
		return
	}

	seqMap := map[string]*ACTGBaseSequence{}
	for _, fileName := range fileNames {
		// load target file
		fileData, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		baseSequences := strings.Split(string(fileData), ">")
		for _, baseSequence := range baseSequences {
			if baseSequence == "" {
				continue
			}
			splitLines := strings.Split(baseSequence, "\n")
			// there should be 3 item after split the last item is empty string
			if len(splitLines) != 3 {
				fmt.Printf("malformed content [%s] in file %s\n", baseSequence, fileName)
				continue
			}
			// remove \r in CRLF mode
			seqName := strings.TrimSuffix(splitLines[0], "\r")
			seqContent := strings.TrimSuffix(splitLines[1], "\r")
			seq, exists := seqMap[seqName]
			if exists {
				seq.Append(seqContent)
			} else {
				seqMap[seqName] = NewACTG(seqName, seqContent)
			}
		}
	}

	osEol := "\r\n"
	mergeContents := ""
	for _, seq := range seqMap {
		mergeContents += fmt.Sprintf(">%s%s%s%s", seq.Name, osEol, seq.Content, osEol)
	}
	ioutil.WriteFile(outputFile, []byte(mergeContents), 0644)
	fmt.Printf("merge date to %s done\n", outputFile)
}
