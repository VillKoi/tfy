package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"tfy/lab2/cmd/solution"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config6.txt", "Used for set path to config file.")
	flag.Parse()

	fmt.Println(configPath)

	dfa, err := PrepareData(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	ПроверкаDFA(dfa)

	dfa = рекурсия_2(dfa)

	fmt.Println(dfa.ABиP[Старт][dfa.КонечноеСостояние][0])
}

type DFA struct {
	НачальноеСостояние string
	КонечноеСостояние  string

	Состояния map[string]string
	Переходы  map[string]struct{}

	КонечныеСостояния map[string]string

	ABиP map[string]map[string][]string
	BAиP map[string]map[string][]string
}

type Equation struct {
	AfterEqual string
}

func PrepareData(configPath string) (DFA, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	data := string(file)

	// убираем первые <>
	data = strings.TrimSuffix(strings.TrimPrefix(data, "<"), ">")
	// сплитим по ,
	//  а мы не можем просто там засплитить, проход по строке?
	// может ли быть мно-во начальных состояний?
	fsep := strings.SplitN(data, ", ", 2)

	dfa := DFA{
		НачальноеСостояние: fsep[0],
	}

	// во втором куске можем засплитить по },{
	ssep := strings.SplitN(fsep[1], "}, {", 2)

	gt := strings.TrimPrefix(ssep[0], "{")

	dfa.Переходы, dfa.Состояния, dfa.ABиP, dfa.BAиP = ПолучениеПереходов(gt)             // {<Q0,a,Q0> <Q0,b,Q1> <Q1,a,Q1> <Q1,b,Q2> <Q2,a,Q2> <Q2,b,Q2>
	dfa.КонечныеСостояния = ПолучениеКонечныхСостояний(strings.TrimSuffix(ssep[1], "}")) // Q1}

	return dfa, nil
}

func ПолучениеПереходов(s string) (переходы map[string]struct{}, состояния map[string]string, ABиP, BAиP map[string]map[string][]string) {
	переходы = make(map[string]struct{})
	состояния = make(map[string]string, 0)
	ABиP = make(map[string]map[string][]string)
	BAиP = make(map[string]map[string][]string)

	названияСостояний := 'A'
	состояния[Старт] = Старт
	состояния[Конец] = Конец

	for strings.Contains(s, "<") {
		first := strings.Index(s, "<")
		second := strings.Index(s, ">")
		переход := s[first+1 : second]
		sep := strings.Split(переход, ",")

		переходы[sep[1]] = struct{}{}

		if _, ok := состояния[sep[0]]; !ok {
			состояния[sep[0]] = string(названияСостояний)
			названияСостояний++
		}

		if _, ok := состояния[sep[2]]; !ok {
			состояния[sep[2]] = string(названияСостояний)
			названияСостояний++
		}

		if _, ok := ABиP[sep[0]]; !ok {
			ABиP[sep[0]] = make(map[string][]string)
		}
		if _, ok := BAиP[sep[2]]; !ok {
			BAиP[sep[2]] = make(map[string][]string)
		}

		ABиP[sep[0]][sep[2]] = append(ABиP[sep[0]][sep[2]], sep[1])
		BAиP[sep[2]][sep[0]] = append(BAиP[sep[2]][sep[0]], sep[1])

		s = strings.TrimPrefix(s, s[:second+1])
	}

	return переходы, состояния, ABиP, BAиP
}

func ПолучениеКонечныхСостояний(s string) map[string]string {
	мапа := make(map[string]string)

	конечныеСостояния := strings.Split(s, ",")

	for i := range конечныеСостояния {
		мапа[конечныеСостояния[i]] = конечныеСостояния[i]
	}

	return мапа
}

func ПроверкаDFA(dfa DFA) {
	printlnDFA(dfa)
	// проверка на отсутствие переходов в состояния из которых не выйти
	for A, B := range dfa.ABиP {
		_, ok := dfa.КонечныеСостояния[A]

		pp := len(B)
		if _, okB := B[A]; okB {
			pp--
		}

		if !ok && pp == 0 {
			fmt.Println("Некорректные входные данные 1: ", A, "из соcтояния нет перехода в конечные состояния", dfa.КонечныеСостояния)
			os.Exit(0)
		}
	}

	// проверка нескольких переходов по одному символу из одного состояния
	for A, B := range dfa.ABиP {
		переходы := make(map[string]struct{})
		for _, pp := range B {
			for i := range pp {
				_, ok := переходы[pp[i]]

				if !ok {
					переходы[pp[i]] = struct{}{}
				} else {
					fmt.Println("Некорректные входные данные 2: ", A, B)
					fmt.Println("---")
					//os.Exit(0)
				}
			}

		}
	}
}

func рекурсия_2(dfa DFA) DFA {
	for A, B := range dfa.ABиP {
		dfa.ABиP[A] = ОбъединениеПереходов(B)
	}

	dfa = добавление_нач_и_кон_состояния(dfa)

	printlnE(dfa)

	// решение номер 2
	dfa = решение_2(dfa)

	for A, B := range dfa.ABиP {
		dfa.ABиP[A] = ОбъединениеПереходов(B)
	}

	return dfa
}

func добавление_нач_и_кон_состояния(dfa DFA) DFA {
	dfa.ABиP[Старт] = map[string][]string{
		dfa.НачальноеСостояние: {"ε"},
	}
	dfa.BAиP[Старт] = map[string][]string{}

	if _, ok := dfa.BAиP[dfa.НачальноеСостояние]; !ok {
		dfa.BAиP[dfa.НачальноеСостояние] = map[string][]string{}
	}

	dfa.BAиP[dfa.НачальноеСостояние][Старт] = []string{"ε"}

	dfa.ABиP[Конец] = map[string][]string{}
	dfa.BAиP[Конец] = map[string][]string{}
	for i := range dfa.КонечныеСостояния {
		if _, ok := dfa.ABиP[dfa.КонечныеСостояния[i]]; !ok {
			dfa.ABиP[dfa.КонечныеСостояния[i]] = make(map[string][]string)
		}
		dfa.ABиP[dfa.КонечныеСостояния[i]][Конец] = []string{"ε"}
		dfa.BAиP[Конец][dfa.КонечныеСостояния[i]] = []string{"ε"}
	}

	dfa.КонечныеСостояния = map[string]string{
		Конец: Конец,
	}

	dfa.НачальноеСостояние = Старт
	dfa.КонечноеСостояние = Конец

	return dfa
}

func ОбъединениеПереходов(BB map[string][]string) map[string][]string {
	for B, pp := range BB {
		newPP := ""
		for i, p := range pp {
			if i == 0 {
				newPP = p
				continue
			}
			newPP = "(" + newPP + " + " + p + ")"
		}
		BB[B] = []string{newPP}
	}
	return BB
}

// формат graphviz
func printlnDFA(dfa DFA) {
	for A, RR := range dfa.ABиP {
		for B, pp := range RR {
			for i := range pp {
				fmt.Println(A, "->", B, "[ label=\"", pp[i], "\" ];")

			}
		}
	}
	fmt.Println("---")
}

func printlnE(dfa DFA) {
	configPath := "../tests/newCon1.txt"

	buffer := bytes.NewBufferString("")

	var strBegin string
	for B, pp := range dfa.ABиP[Старт] {
		for i := range pp {
			if strBegin == "" {
				strBegin = dfa.Состояния[Старт] + " = " + pp[i] + dfa.Состояния[B]
				continue
			}

			strBegin += " + " + pp[i] + dfa.Состояния[B]
		}
	}
	fmt.Println(strBegin)
	buffer.WriteString(strBegin + "\n")

	for A := range dfa.Состояния {
		if A == Старт {
			continue
		}

		var str string

		if len(dfa.ABиP[A]) == 0 {
			str = dfa.Состояния[A] + " = " + "ε"
		}

		for B, pp := range dfa.ABиP[A] {
			for i := range pp {
				if str == "" {
					str = dfa.Состояния[A] + " = " + pp[i] + dfa.Состояния[B]
					continue
				}

				str += " + " + pp[i] + dfa.Состояния[B]
			}
		}
		fmt.Println(str)
		buffer.WriteString(str + "\n")

	}
	fmt.Println("---")

	err := ioutil.WriteFile(configPath, buffer.Bytes(), 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	equationSystem, err := solution.GetData(configPath)
	if err != nil {
		fmt.Println("Решение с помощью слау", err)
		return
	}

	ответ := solution.Pешение(equationSystem)

	fmt.Println(ответ.Var + "=" + ответ.Regex)
}

var Старт = "X"
var Конец = "Z"

func решение_2(dfa DFA) DFA {
	for len(dfa.ABиP) > 2 {
		// берем вершину в середине
		for R, BB := range dfa.ABиP {
			for C, D := range dfa.ABиP {
				dfa.ABиP[C] = ОбъединениеПереходов(D)
			}
			for D, C := range dfa.BAиP {
				dfa.BAиP[D] = ОбъединениеПереходов(C)
			}

			printlnDFA(dfa)
			// все переходы в эту вершину
			AA := dfa.BAиP[R]
			if len(BB) == 0 || len(AA) == 0 {
				continue
			}

			середина := ""
			if ppR, ok := dfa.ABиP[R][R]; ok {
				if strings.HasPrefix(ppR[0], "(") {
					середина = ppR[0] + "*"
				} else {
					середина = "(" + ppR[0] + ")*"
				}
			}

			i := 0
			for A, ppA := range AA {
				if A == R {
					continue
				}
				for B, ppB := range BB {
					if B == R {
						continue
					}
					i++

					newP := ""

					if ppA[0] == "ε" {
						newP = середина + ppB[0]
					} else if ppB[0] == "ε" {
						newP = ppA[0] + середина
					} else {
						newP = ppA[0] + середина + ppB[0]
					}

					if _, ok := dfa.ABиP[A]; !ok {
						dfa.ABиP[A] = map[string][]string{}
					}
					if _, ok := dfa.BAиP[B]; !ok {
						dfa.BAиP[B] = map[string][]string{}
					}

					dfa.ABиP[A][B] = append(dfa.ABиP[A][B], newP)
					dfa.BAиP[B][A] = append(dfa.BAиP[B][A], newP)
				}
			}
			// прошли все R - удаляем
			if i > 0 {
				delete(dfa.ABиP, R)
				delete(dfa.BAиP, R)

				for A := range AA {
					delete(dfa.ABиP[A], R)
				}

				for B := range BB {
					delete(dfa.BAиP[B], R)
				}
			}
		}
	}

	return dfa
}
