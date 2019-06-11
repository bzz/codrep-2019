package main

// Print formatting errors in every file in 'learning' dataset

import (
	"encoding/json"
	"unicode/utf8"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	offsetFilename = "out.txt"
	contextChars   = 15
)

var (
	dataset = filepath.Join("Datasets", "learning") // TODO(bzz): add [-d] or [-datasets=]
	histogram = make(map[string]int)
)

func main() {
	offsetFile := filepath.Join(dataset, offsetFilename)
	offsets, err := ioutil.ReadFile(offsetFile)
	if err != nil {
		log.Fatalf("failed reading %s, err: %s", offsetFile, err)
	}

	for lineNumber, fmtError := range strings.Split(string(offsets), "\n") {
		if fmtError == "" {
			continue
		}

		filename := fmt.Sprintf("%d.txt", lineNumber)
		errorOffset, err := strconv.Atoi(fmtError)
		if err != nil {
			log.Fatalf("file %s does not have 1 offset (%s), err: %s\n", filename, fmtError, err)
		}

		data, err := ioutil.ReadFile(filepath.Join(dataset, filename))
		if err != nil {
			log.Printf("faild to read %s, err: %s\n", filename, err)
		}

		dataStr := string(data)
		data = nil
		if !utf8.ValidString(dataStr) {
			log.Fatalf("file is not valid utf8 %s\n", filename)
		}


		if len(dataStr) < errorOffset {
			log.Fatalf("offset is bigger then file - file:%s len:%d offset:%d", filename, len(dataStr), errorOffset)
		}

		if lineNumber%3 == 0 && lineNumber!=0 {
			fmt.Println()
		}
		printLine(dataStr, lineNumber, errorOffset)

	}
	fmt.Println()
	b, _ := json.MarshalIndent(histogram, "", "  ")
	fmt.Printf("%d different errors\n %s\n", len(histogram), b)
}

func printLine(str string, lineNumber, errorOffset int) {
	var br, er, ar []rune
	index := 0
	for _, runeVal := range str {
		if max(0, errorOffset-contextChars) <= index && index < errorOffset-1 {
			br = append(br, runeVal) //str[max(0, errorOffset-contextChars) : errorOffset-1]
		} else if errorOffset-1 <= index && index < errorOffset {
			er = append(er, runeVal) // str[errorOffset-1 : errorOffset]
		} else if errorOffset <= index && index < min(errorOffset+contextChars, len(str)) {
			ar = append(ar, runeVal) // str[errorOffset:min(errorOffset+contextChars, len(str))]
		}
		index++
	}
	before := keepNewlines(br)
	fmt.Printf("%6d %17s", lineNumber, before)

	fmtError := keepNewlines(er)
	fmt.Printf("->%2s<-", fmtError)
	// histogram[fmtError]++
	if fmtError == `\n` || fmtError == `\t` || fmtError == " " {
		histogram[fmtError]++
	} else {
		histogram["Ø"]++
	}

	after := keepNewlines(ar)
	fmt.Printf("%-19s\"", after)
}

func keepNewlines(sr []rune) string {
	st := strings.Replace(string(sr), "\n", `\n`, -1)
	return strings.Replace(st, "\t", `\t`, -1)
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
