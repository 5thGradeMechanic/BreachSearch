package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"
)

var patterns string
var inputs string

func init() {
	flag.StringVar(&patterns, "p", "None", "Please provide absolute file path to patterns file")
	flag.StringVar(&inputs, "i", "None", "Please provide absolute path to inputs directory")
}

func usage() {
	//Display usage when someone doesn't pass a flag
	fmt.Printf("Usage: %s [OPTIONS] arguments ....\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)

}

func getRegex() *regexp.Regexp {
	//Compile regex and return it
	var re = *regexp.MustCompile(`@(?:[a-zA-Z0-9-]+\.)?([A-Za-z0-9'.]){2,}`)
	return &re
}

func binSearch(words []string, word string) int {
	var lo int = 0
	var hi int = len(words) - 1

	for lo <= hi {
		var mid int = lo + (hi-lo)/2
		var midValue string = words[mid]

		if compare(midValue, word) == 0 {
			return mid
		} else if compare(midValue, word) > 0 {
			// We want to use the left half of our list
			hi = mid - 1
		} else {
			// We want to use the right half of our list
			lo = mid + 1
		}
	}

	// If we get here we tried to look at an invalid sub-list
	// which means the number isn't in our list.
	return -1
}

func compare(a, b string) int {
	var aLow string = strings.ToLower(a)
	var bLow string = strings.ToLower(b)
	if aLow == bLow {
		return 0
	} else if aLow < bLow {
		return -1
	} else {
		return 1
	}
}

func StripWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func readInputsDirectory(dirname string) ([]string, error) {
	//Read in all the filenames in the Inputs directory
	var inputs []string
	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		inputs = append(inputs, file.Name())
	}
	return inputs, err
}

func readPatterns(path string) ([]string, error) {
	//Read in our patterns that will be used to look for matches
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		patterns = append(patterns, StripWhiteSpace(scanner.Text()))
	}
	return patterns, scanner.Err()
}

func outputFile(results []string, path string) error {
	//Dump our slice out to the file
	file, err := os.Create(("results_" + path))
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, result := range results {
		fmt.Fprintln(w, result)
	}
	return w.Flush()
}

func checkPatterns(path string, patterns []string, input string, wg *sync.WaitGroup) ([]string, error) {
	//The nested loop goes over file checking if the pattern is a substring of the line
	defer wg.Done()
	re := getRegex()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var results []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			var index int
			index = binSearch(patterns, re.FindString(scanner.Text()))
			if index < 0 {
				//fmt.Println("There was no email in ", scanner.Text(), "could not be found!")
			} else {
				//fmt.Println("Email match in string", scanner.Text(), "at index:", index, patterns[index])
				resultant := (input + ";/;" + patterns[index] + ";/;" + scanner.Text())
				results = append(results, resultant)
			}
		}

	}
	outputFile(results, input)
	return results, scanner.Err()
}

func main() {

	//Parse flags and setup CLI params
	flag.Parse()
	if inputs == "None" || patterns == "None" {
		usage()
	}

	//Read in our patterns from file
	var pattern_file string = patterns

	patterns, err := readPatterns(pattern_file)
	if err != nil {
		log.Fatalf("readLines :%s", err)
	}
	//Remove duplicates
	patterns = removeDuplicates(patterns)

	sort.Strings(patterns)

	//Get all the input files
	var inputs_directory string = inputs
	inputs, err := readInputsDirectory(inputs_directory)
	if err != nil {
		log.Fatalf("Inputs failed :%s", err)
	}

	//Setup waitgroup for synchronized processing
	var wg sync.WaitGroup
	wg.Add(len(inputs))

	//Concurrently search all inputs for pattern matches, will use multiple cores
	for _, input := range inputs {
		go checkPatterns((inputs_directory + input), patterns, input, &wg)
	}

	wg.Wait()
}
