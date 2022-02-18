package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config6.txt", "Used for set path to config file.")
	flag.Parse()

	_, _, unifier := unification(configPath)

	unifierString := printTerm(unifier)

	fmt.Println("Unifier: ", unifierString)
}

func printTerm(unifier Term) string {
	if len(unifier.innerterm) == 0 {
		return unifier.constructor
	}

	unifierString := unifier.constructor + "("

	for i, v := range unifier.innerterm {
		if i == len(unifier.innerterm)-1 {
			unifierString += printTerm(v) + ")"
			continue
		}

		unifierString += printTerm(v) + ","
	}

	return unifierString
}

func unification(configPath string) (changeTerm1, changeTerm2 []string, term Term) {
	t, err := Get(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	term1 := t.startParce(t.First)
	term2 := t.startParce(t.Second)

	unifier := t.unif(term1, term2)

	return t.changeTerm1, t.changeTerm2, unifier
}

// easytags main.go json
type TRS struct {
	Constructors map[string]int      `json:"constructors"`
	Variables    map[string]struct{} `json:"variables"`
	First        string              `json:"first"`
	Second       string              `json:"second"`
	Term         string              `json:"term"`

	changeTerm1 []string
	changeTerm2 []string
}

func Get(configPath string) (TRS, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	data := strings.Split(string(file), "\n")

	constructors := map[string]int{}
	for _, v := range strings.Split(strings.ReplaceAll(strings.TrimPrefix(data[0], "constructors ="), " ", ""), ",") {
		parts := strings.SplitN(v, "(", 2)
		signature, err := strconv.Atoi(strings.TrimSuffix(parts[1], ")"))
		if err != nil {
			log.Panic(err)
		}

		constructors[parts[0]] = signature
	}

	variables := map[string]struct{}{}
	for _, v := range strings.Split(strings.ReplaceAll(strings.TrimPrefix(data[1], "variables ="), " ", ""), ",") {
		variables[v] = struct{}{}
	}

	t := TRS{
		Constructors: constructors,
		Variables:    variables,
		First:        strings.ReplaceAll(strings.TrimSpace(strings.TrimPrefix(data[2], "first =")), " ", ""),
		Second:       strings.ReplaceAll(strings.TrimSpace(strings.TrimPrefix(data[3], "second =")), " ", ""),
	}

	return t, nil
}

type Term struct {
	Fullterm    string
	constructor string
	signature   string
	innerterm   []Term
}

func (t TRS) startParce(termS string) Term {
	term := Term{
		Fullterm: termS,
	}

	return t.parceTree(term)
}

func (t TRS) parceTree(term Term) Term {
	if !strings.Contains(term.Fullterm, "(") {
		term.constructor = term.Fullterm
		return term
	}

	parts := strings.SplitN(term.Fullterm, "(", 2)

	term.constructor = parts[0]

	if strings.EqualFold(parts[1], ")") {
		return term
	}

	term.signature = strings.TrimSuffix(parts[1], ")")

	fmt.Println("constructor: ", term.constructor, "\n",
		"signature: ", term.signature)

	lists := []string{}
	if !strings.Contains(term.signature, "(") {
		lists = strings.Split(term.signature, ",")
	} else {
		bracketCounter := 0
		slice := []rune(term.signature)
		lastCeparate := 0
		for i := range slice {
			if i == len(slice)-1 {
				lists = append(lists, string(slice[lastCeparate:]))
				break
			}

			if strings.EqualFold(string(slice[i]), "(") {
				bracketCounter++
			} else if strings.EqualFold(string(slice[i]), ")") {
				bracketCounter--
			} else if bracketCounter == 0 && strings.EqualFold(string(slice[i]), ",") {
				lists = append(lists, string(slice[lastCeparate:i]))
				lastCeparate = i + 1
			}
		}
	}

	if len(lists) != t.Constructors[term.constructor] {
		fmt.Println("invalid term", term.Fullterm)
		os.Exit(0)
	}

	for i := range lists {
		term.innerterm = append(term.innerterm, t.parceTree(Term{
			Fullterm: lists[i],
		}))
	}

	return term
}

const (
	constructor = iota
	variable
	constant
)

func (trs TRS) whichType(v string) int {
	if number, ok := trs.Constructors[v]; ok {
		if number == 0 {
			return constant
		}
		return constructor
	} else {
		if _, ok := trs.Variables[v]; ok {
			return variable
		}
	}

	return constant
}

func (t TRS) unif(term1, term2 Term) Term {
	if variable == t.whichType(term1.constructor) {
		if term1.constructor != term2.constructor {
			t.changeTerm1 = append(t.changeTerm1, term1.Fullterm+":="+term2.Fullterm)

			fmt.Println(term1.Fullterm + ":=" + term2.Fullterm + ", ")
		}

		return term2
	}

	if variable == t.whichType(term2.constructor) {
		if term1.constructor != term2.constructor {
			t.changeTerm2 = append(t.changeTerm2, term2.Fullterm+":="+term1.Fullterm)

			fmt.Println(term2.Fullterm + ":=" + term1.Fullterm + ", ")
		}

		return term1
	}

	if constructor == t.whichType(term1.constructor) && constructor == t.whichType(term2.constructor) &&
		term1.constructor == term2.constructor {
		args := []Term{}

		for i := range term1.innerterm {
			args = append(args, t.unif(term1.innerterm[i], term2.innerterm[i]))
		}

		return Term{
			Fullterm:    term1.Fullterm,
			constructor: term1.constructor,
			signature:   term1.signature,
			innerterm:   args,
		}
	}

	if constant == t.whichType(term1.constructor) && constant == t.whichType(term2.constructor) &&
		term1.signature == term2.signature {
		return term1
	}

	log.Fatal("Unable to unify")
	return Term{}
}
