package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	var dfa1Path string
	var dfa2Path string

	flag.StringVar(&dfa1Path, "c1", "DFA_3_1.txt", "Used for set path to config file.")
	flag.StringVar(&dfa2Path, "c2", "DFA_3_2.txt", "Used for set path to config file.")
	flag.Parse()

	fmt.Println(dfa1Path)
	fmt.Println(dfa2Path)

	dfa_1, err := PrepareDFA(dfa1Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	dfa_2, err := PrepareDFA(dfa2Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	dfa_1 = достройкаЛовушек(dfa_1)

	printlnDFA(dfa_1)

	dfa_2 = достройкаЛовушек(dfa_2)

	printlnDFA(dfa_2)

	// printlnDFA(пересечение_DFA(dfa_1, dfa_2))

	d_dfa_1_ := дополнение(dfa_1)

	printlnDFA(d_dfa_1_)

	d_dfa_2_ := дополнение(dfa_2)

	printlnDFA(d_dfa_2_)

	p_d_1_d_d_2 := пересечение_DFA(dfa_1, d_dfa_2_)

	p_d_2_d_d_1 := пересечение_DFA(dfa_2, d_dfa_1_)

	printlnDFA(p_d_1_d_d_2)

	printlnDFA(p_d_2_d_d_1)

	пер_1 := проверкаДостижимостиВсехКонечных(p_d_1_d_d_2)
	пер_2 := проверкаДостижимостиВсехКонечных(p_d_2_d_d_1)

	if !пер_1 && !пер_2 {
		fmt.Println("языки эквивалентны")
	}

}

type DFA struct {
	НачальноеСостояние string
	КонечноеСостояние  string

	Состояния map[string]string

	ПарыСостояний map[string][2]string

	// алфавит
	Переходы map[string]struct{}

	КонечныеСостояния map[string]string

	ABиP map[string]map[string][]string
	BAиP map[string]map[string][]string
	APиB map[string]map[string][]string
}

func PrepareDFA(path string) (DFA, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return DFA{}, err
	}

	data := string(file)

	fsep := strings.SplitN(data, ", ", 2)

	dfa := DFA{
		НачальноеСостояние: strings.TrimPrefix(fsep[0], "<"),
	}

	ssep := strings.SplitN(fsep[1], "}, {", 2)

	gt := strings.TrimPrefix(ssep[0], "{")

	dfa = ПолучениеПереходов(gt, dfa)                                                     // {<Q0,a,Q0> <Q0,b,Q1> <Q1,a,Q1> <Q1,b,Q2> <Q2,a,Q2> <Q2,b,Q2>
	dfa.КонечныеСостояния = ПолучениеКонечныхСостояний(strings.TrimSuffix(ssep[1], "}>")) // Q1}

	return dfa, nil
}

func ПолучениеПереходов(s string, dfa DFA) DFA {
	переходы := make(map[string]struct{})
	состояния := make(map[string]string, 0)
	ABиP := make(map[string]map[string][]string)
	BAиP := make(map[string]map[string][]string)

	APиB := make(map[string]map[string][]string)

	s = strings.ReplaceAll(s, " ", "")

	for strings.Contains(s, "<") {
		first := strings.Index(s, "<")

		if !strings.HasPrefix(s, "<") {
			fmt.Println("Что-то лишнее между > <:", s[:first])
			os.Exit(0)
		}

		second := strings.Index(s, ">")
		переход := s[first+1 : second]
		sep := strings.Split(переход, ",")

		переходы[sep[1]] = struct{}{}

		if _, ok := состояния[sep[0]]; !ok {
			состояния[sep[0]] = sep[0]
		}

		if _, ok := состояния[sep[2]]; !ok {
			состояния[sep[2]] = sep[2]
		}

		if _, ok := ABиP[sep[0]]; !ok {
			ABиP[sep[0]] = make(map[string][]string)
		}
		if _, ok := BAиP[sep[2]]; !ok {
			BAиP[sep[2]] = make(map[string][]string)
		}

		if _, ok := APиB[sep[0]]; !ok {
			APиB[sep[0]] = make(map[string][]string)
		}

		ABиP[sep[0]][sep[2]] = append(ABиP[sep[0]][sep[2]], sep[1])
		BAиP[sep[2]][sep[0]] = append(BAиP[sep[2]][sep[0]], sep[1])

		APиB[sep[0]][sep[1]] = append(ABиP[sep[0]][sep[1]], sep[2])

		s = strings.TrimPrefix(s, s[:second+1])
	}

	dfa.Переходы = переходы
	dfa.Состояния = состояния
	dfa.ABиP = ABиP
	dfa.BAиP = BAиP
	dfa.APиB = APиB

	return dfa
}

func ПолучениеКонечныхСостояний(s string) map[string]string {
	мапа := make(map[string]string)

	s = strings.ReplaceAll(s, " ", "")
	конечныеСостояния := strings.Split(s, ",")

	for i := range конечныеСостояния {
		мапа[конечныеСостояния[i]] = конечныеСостояния[i]
	}

	return мапа
}

func ПроверкаDFA(dfa DFA) {
	printlnDFA(dfa)

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

func достройкаЛовушек(dfa DFA) DFA {
	nameЛовушка := "Л"
	iЛовушки := 0

	ловушка := ""

	for state := range dfa.Состояния {
		newЛовушка := true
		for p := range dfa.Переходы {
			_, ok := dfa.APиB[state][p]

			if !ok {
				//	создаем ловушку
				if newЛовушка {
					ловушка = nameЛовушка + strconv.Itoa(iЛовушки)
					iЛовушки++
					newЛовушка = false

					dfa.ABиP[ловушка] = map[string][]string{}
					dfa.BAиP[ловушка] = map[string][]string{}
					dfa.APиB[ловушка] = map[string][]string{}

					for pp := range dfa.Переходы {
						dfa.ABиP[ловушка][ловушка] = append(dfa.ABиP[ловушка][ловушка], pp)
						dfa.BAиP[ловушка][ловушка] = append(dfa.BAиP[ловушка][ловушка], pp)
						dfa.APиB[ловушка][pp] = append(dfa.APиB[ловушка][pp], ловушка)
					}
					dfa.Состояния[ловушка] = ловушка
				}

				if _, ok := dfa.ABиP[state]; !ok {
					dfa.ABиP[state] = map[string][]string{}
				}

				if _, ok := dfa.APиB[state]; !ok {
					dfa.APиB[state] = map[string][]string{}
				}

				dfa.ABиP[state][ловушка] = append(dfa.ABиP[state][ловушка], p)
				dfa.BAиP[ловушка][p] = append(dfa.BAиP[ловушка][p], state)
				dfa.APиB[state][p] = append(dfa.APиB[state][p], ловушка)
			}
		}
	}

	return dfa
}

// все конечные -> неконечные, неконечные -> конечные
func дополнение(dfa DFA) DFA {
	d_dfa := DFA{
		НачальноеСостояние: dfa.НачальноеСостояние,
		КонечныеСостояния:  make(map[string]string),

		Состояния: make(map[string]string, 0),
		Переходы:  make(map[string]struct{}),
		ABиP:      make(map[string]map[string][]string),
		BAиP:      make(map[string]map[string][]string),
		APиB:      make(map[string]map[string][]string),
	}

	for k, v := range dfa.ABиP {
		d_dfa.ABиP[k] = v
	}

	for k, v := range dfa.BAиP {
		d_dfa.BAиP[k] = v
	}

	for k, v := range dfa.APиB {
		d_dfa.APиB[k] = v
	}

	for k := range dfa.Переходы {
		d_dfa.Переходы[k] = struct{}{}
	}

	for k, v := range dfa.Состояния {
		d_dfa.Состояния[k] = v

		if _, ok := dfa.КонечныеСостояния[k]; !ok {
			d_dfa.КонечныеСостояния[k] = v
		}
	}

	return d_dfa
}

func пересечение_DFA(dfa_1, dfa_2 DFA) DFA {
	начальноеСостояние := dfa_1.НачальноеСостояние + dfa_2.НачальноеСостояние

	new_dfa := DFA{
		НачальноеСостояние: начальноеСостояние,

		КонечныеСостояния: map[string]string{},

		Состояния: map[string]string{
			начальноеСостояние: начальноеСостояние,
		},

		ПарыСостояний: map[string][2]string{
			начальноеСостояние: [2]string{
				dfa_1.НачальноеСостояние, dfa_2.НачальноеСостояние,
			},
		},

		ABиP: map[string]map[string][]string{},
		BAиP: map[string]map[string][]string{},
		APиB: map[string]map[string][]string{},
	}

	m := -1

	все_переходы := dfa_1.Переходы

	for k, v := range dfa_2.Переходы {
		все_переходы[k] = v
	}

	for m != len(new_dfa.Состояния) {
		m = len(new_dfa.Состояния)

		printlnDFA(new_dfa)

		for state := range new_dfa.Состояния {
			for p := range все_переходы {
				_, ok := new_dfa.APиB[state][p]

				if !ok {
					states := new_dfa.ПарыСостояний[state]
					state_1, state_2 := states[0], states[1]

					BB_1, ok_1 := dfa_1.APиB[state_1][p]

					BB_2, ok_2 := dfa_2.APиB[state_2][p]

					if !ok_1 || !ok_2 {
						continue
					}

					new_state_B := BB_1[0] + BB_2[0]

					if _, ok := new_dfa.ABиP[state]; !ok {
						new_dfa.ABиP[state] = map[string][]string{}
					}

					if _, ok := new_dfa.APиB[state]; !ok {
						new_dfa.APиB[state] = map[string][]string{}
					}

					if _, ok := new_dfa.BAиP[new_state_B]; !ok {
						new_dfa.BAиP[new_state_B] = map[string][]string{}
					}

					new_dfa.Состояния[new_state_B] = new_state_B
					new_dfa.ПарыСостояний[new_state_B] = [2]string{BB_1[0], BB_2[0]}

					new_dfa.ABиP[state][new_state_B] = append(new_dfa.ABиP[state][new_state_B], p)
					new_dfa.BAиP[new_state_B][p] = append(new_dfa.BAиP[new_state_B][p], state)
					new_dfa.APиB[state][p] = append(new_dfa.APиB[state][p], new_state_B)

					_, ok_1 = dfa_1.КонечныеСостояния[state_1]
					_, ok_2 = dfa_2.КонечныеСостояния[state_2]

					if ok_1 && ok_2 {
						new_dfa.КонечныеСостояния[state] = state
					}

					_, ok_1 = dfa_1.КонечныеСостояния[BB_1[0]]
					_, ok_2 = dfa_2.КонечныеСостояния[BB_2[0]]

					if ok_1 && ok_2 {
						new_dfa.КонечныеСостояния[new_state_B] = new_state_B
					}
				}
			}
		}
	}

	printlnDFA(new_dfa)

	return new_dfa
}

func проверкаДостижимостиВсехКонечных(dfa DFA) bool {
	if len(dfa.КонечноеСостояние) == 0 {
		return false
	}

	// пусто ->	все конечные не достижимы

	достижимые := true

	for k := range dfa.КонечныеСостояния {
		mapД := map[string]struct{}{
			k: {},
		}
		m := -1

		for m != len(mapД) {
			m = len(mapД)

			for д := range mapД {
				for A := range dfa.BAиP[д] {
					mapД[A] = struct{}{}
				}
			}
		}

		if _, ok := mapД[dfa.НачальноеСостояние]; !ok {
			return false
		}
	}

	return достижимые
}

func проверкаНаПустотуПересечения(dfa DFA) bool {
	dfa = удалениеНедостижимыхИНепорождающих(dfa)

	printlnDFA(dfa)

	//  все конечные состояния недостижимы

	if len(dfa.Состояния) == 0 {
		return true
	}

	return false
}

func удалениеНедостижимыхИНепорождающих(dfa DFA) DFA {
	m := -1

	for m != len(dfa.Состояния) {
		m = len(dfa.Состояния)

		dfa = получениеДостижимыхИУдалениеНедостижимых(dfa)

		dfa = получениеПорождающих(dfa)
	}

	return dfa
}

func получениеДостижимыхИУдалениеНедостижимых(dfa DFA) DFA {
	достижимые_нетермы := map[string]struct{}{
		dfa.НачальноеСостояние: {},
	}

	m := -1
	for m != len(достижимые_нетермы) {
		m = len(достижимые_нетермы)

		for D := range достижимые_нетермы {
			for B := range dfa.ABиP[D] {
				достижимые_нетермы[B] = struct{}{}
			}
		}
	}
	недостижимые := map[string]struct{}{}

	for A := range dfa.Состояния {
		if _, ok := достижимые_нетермы[A]; !ok {
			недостижимые[A] = struct{}{}
		}
	}

	for B := range недостижимые {
		AA := dfa.BAиP[B]

		for A := range AA {
			delete(dfa.ABиP[A], B)

			for p, BB := range dfa.APиB[A] {
				newBB := []string{}
				for _, v := range BB {
					if v == B {
						continue
					}
					newBB = append(newBB, v)
				}
				dfa.APиB[A][p] = newBB
			}
		}
		delete(dfa.BAиP, B)
		delete(dfa.КонечныеСостояния, B)
		delete(dfa.Состояния, B)
	}
	return dfa
}

func получениеПорождающих(dfa DFA) DFA {
	m := -1
	непорождающие := map[string]struct{}{}

	for m != len(непорождающие) {
		m = len(непорождающие)
		// нахождение
		for B, AA := range dfa.BAиP {
			_, ok := dfa.КонечныеСостояния[B]
			if ok {
				continue
			}

			неп := true
			for A := range AA {
				if B != A {
					неп = false
				}
			}
			if неп && len(dfa.BAиP[B]) == 0 {
				непорождающие[B] = struct{}{}
			}
		}
		// удаление
		for B := range непорождающие {
			AA := dfa.BAиP[B]

			for A := range AA {
				delete(dfa.ABиP[A], B)

				for p, BB := range dfa.APиB[A] {
					newBB := []string{}
					for _, v := range BB {
						if v == B {
							continue
						}

						newBB = append(newBB, v)
					}

					dfa.APиB[A][p] = newBB
				}
			}

			delete(dfa.BAиP, B)
			delete(dfa.КонечныеСостояния, B)
			delete(dfa.Состояния, B)
		}
	}
	return dfa
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

	for state := range dfa.КонечныеСостояния {
		fmt.Println(state, "[shape=Msquare];")
	}

	fmt.Println(dfa.НачальноеСостояние, "[shape=Mdiamond];")

	fmt.Println("---")
}

// func PrepareDataRegex(path string) error {
// file, err := ioutil.ReadFile(path)
// if err != nil {
// 	return err
// }

// data := strings.Split(string(file), "\n")

// 	if len(data) != 4 {
// 		return errors.New("Некорректные входные данные")
// 	}

// 	first_regex := data[1]
// 	second_regex := data[3]

// 	for i := range data {
// 		all_string := strings.ReplaceAll(data[i], " ", "")

// 		if strings.EqualFold(all_string, "") {
// 			continue
// 		}

// 	}

// 	return nil
// }
