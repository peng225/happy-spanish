package main

import (
	"fmt"
	"io"
	"os"

	"sigs.k8s.io/yaml"
)

type Word struct {
	Original          string
	InJapanese        string
	Present           []string
	PastParticiple    string
	DottedPast        []string
	LinearPast        []string
	PresentParticiple string
	Future            []string
}

var (
	idToSubject []string
)

func init() {
	idToSubject = []string{
		"yo",
		"tú",
		"él/ella/usted",
		"nosotoros/as",
		"vosotoros/as",
		"ellos/ellas/ustedes",
	}
}

func printTOC(f io.Writer, words []*Word) {
	fmt.Fprintf(f, "## 目次\n")
	for _, word := range words {
		fmt.Fprintf(f, "- [%s](#%s)\n", word.Original, word.Original)
	}
}

func printVerbs(f io.Writer, tense string, verbs []string) {
	fmt.Fprintf(f, "- %s\n", tense)
	printedAny := false
	for i, verb := range verbs {
		if verb == "" {
			continue
		}
		fmt.Fprintf(f, "  - (%s) %s\n", idToSubject[i], verb)
		printedAny = true
	}
	if !printedAny {
		fmt.Fprintln(f, "  - (skip)")
	}
}

func printOneVerb(f io.Writer, tense string, verb string) {
	fmt.Fprintf(f, "- %s\n", tense)
	if verb == "" {
		fmt.Fprintln(f, "  - (skip)")
		return
	}
	fmt.Fprintf(f, "  - %s\n", verb)
}

func printAsMarkDown(words []*Word, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	fmt.Fprintln(f, "# 不規則動詞の活用")
	printTOC(f, words)
	for _, word := range words {
		fmt.Fprintf(f, "## %s\n", word.Original)
		fmt.Fprintf(f, "- 意味: %s\n", word.InJapanese)

		printVerbs(f, "現在形", word.Present)
		printOneVerb(f, "過去分詞", word.PastParticiple)
		printVerbs(f, "点過去", word.DottedPast)
		printVerbs(f, "線過去", word.LinearPast)
		printOneVerb(f, "現在分詞", word.PresentParticiple)
		printVerbs(f, "未来系", word.Future)
	}
	return nil
}

func main() {
	data, err := os.ReadFile("dict.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var words []*Word
	if err := yaml.Unmarshal(data, &words); err != nil {
		panic(fmt.Errorf("failed to parse yaml: %w", err))
	}

	err = printAsMarkDown(words, "verbs.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
