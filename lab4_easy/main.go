package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	var cfgPath string
	var rgPath string

	flag.StringVar(&cfgPath, "c", "CFG_4.txt", "Used for set path to config file.")
	flag.StringVar(&rgPath, "r", "RG_4.txt", "Used for set path to config file.")
	flag.Parse()

	fmt.Println(cfgPath)
	fmt.Println(rgPath)

	cfg, err := PrepareDataCFG(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	rg, err := PrepareDataRG(rgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fsa := RGtoFSA(rg)

	новые_порождающие := базовый_алгоритм(cfg, fsa)

	fmt.Println("----")
	fmt.Println("Ответ")
	for k := range новые_порождающие {
		print(новые_порождающие[k])
	}
}

func print(kkk []RK) {
	for i := range kkk {
		fmt.Print(kkk[i].In, " -> ")
		for j := range kkk[i].Out {
			fmt.Print(kkk[i].Out[j])
		}
		fmt.Println()
	}
}

type CFG struct {
	StarmNTerm string

	Rules []Rule

	map_Rules map[string]map[string][]string

	NTerm map[string]struct{}
}
type Rule struct {
	Nterm string // in
	Out   string

	terms  []string
	Nterms []string
}

var re_cfg_term = regexp.MustCompile(`[a-z]`)
var re_cfg_Nterm = regexp.MustCompile(`[A-Z][0-9]*`)

func PrepareDataCFG(path string) (CFG, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return CFG{}, err
	}

	data := strings.Split(string(file), "\n")

	cfg := CFG{
		NTerm: map[string]struct{}{},
	}

	for i := range data {
		all_string := strings.ReplaceAll(data[i], " ", "")

		if strings.EqualFold(all_string, "") {
			continue
		}
		spl := strings.Split(all_string, "->")

		cfg.NTerm[spl[0]] = struct{}{}

		if cfg.StarmNTerm == "" {
			cfg.StarmNTerm = spl[0]
		}

		terms := re_cfg_term.FindAllString(spl[1], -1)

		Nterms := re_cfg_Nterm.FindAllString(spl[1], -1)

		cfg.Rules = append(cfg.Rules, Rule{
			Nterm: spl[0],
			Out:   spl[1],

			terms:  terms,
			Nterms: Nterms,
		})
	}

	return cfg, nil
}

type RG struct {
	Rules  []Rule
	Nterm  map[string]struct{}
	Letter map[string]struct{}

	FirstState string
}

func PrepareDataRG(path string) (RG, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return RG{}, err
	}

	data := strings.Split(string(file), "\n")

	rg := RG{
		Nterm:  map[string]struct{}{},
		Letter: map[string]struct{}{},
	}

	for i := range data {
		all_string := strings.ReplaceAll(data[i], " ", "")

		if strings.EqualFold(all_string, "") {
			continue
		}
		spl := strings.Split(all_string, "->")

		rg.Nterm[spl[0]] = struct{}{}

		rg.Rules = append(rg.Rules, Rule{
			Nterm: spl[0],
			Out:   spl[1],
		})

		if rg.FirstState == "" {
			rg.FirstState = spl[0]
		}
	}

	return rg, nil
}

type FSA struct {
	FirstState string
	// если по состоянию нет перехода оно конечное
	EndStates []string

	Состояния map[string]string

	ABиP map[string]map[string][]string
	BAиP map[string]map[string][]string

	PиAB map[string]map[string][]string
}

var (
	firslState = "FIRST"
	endState   = "END"
)

var re_rg_letter = regexp.MustCompile(`[a-z]`)
var re_rg_Nterm = regexp.MustCompile(`[A-Z]`)

func RGtoFSA(rg RG) FSA {
	fsa := FSA{
		FirstState: rg.FirstState,

		Состояния: map[string]string{},

		ABиP: map[string]map[string][]string{},
		BAиP: map[string]map[string][]string{},

		PиAB: map[string]map[string][]string{},
	}

	for _, rule := range rg.Rules {
		fsa.Состояния[rule.Nterm] = rule.Nterm

		B := re_rg_Nterm.FindString(rule.Out)
		if B == "" {
			B = endState
		}
		fsa.Состояния[B] = B

		pp := re_rg_letter.FindString(rule.Out)

		if _, ok := fsa.ABиP[rule.Nterm]; !ok {
			fsa.ABиP[rule.Nterm] = map[string][]string{}
		}

		if _, ok := fsa.PиAB[pp]; !ok {
			fsa.PиAB[pp] = map[string][]string{}
		}

		fsa.ABиP[rule.Nterm][B] = append(fsa.ABиP[rule.Nterm][B], pp)

		fsa.PиAB[pp][rule.Nterm] = append(fsa.PиAB[pp][rule.Nterm], B)
	}

	for state := range fsa.Состояния {
		pp, ok_есть_переходы := fsa.ABиP[state]
		var ok_переход_по_себе bool

		if ok_есть_переходы {
			_, ok_переход_по_себе = pp[state]
		}

		if !ok_есть_переходы || len(pp) == 1 && ok_переход_по_себе {
			fsa.EndStates = append(fsa.EndStates, state)
		}

	}

	return fsa
}

func базовый_алгоритм(cfg CFG, fsa FSA) map[Комбинация][]RK {
	// step 1: X -> t
	комбинации_map_slice, rules_2 := получениеПростыхПравил(cfg, fsa)

	комбинации_2 := получениеОстальныхПравил(fsa, rules_2)

	// step 2: удаление неподождающих
	порождающие := получениеПорождающих(комбинации_map_slice, комбинации_2)

	// сборка всех достижимых
	достижимые_2 := map[Комбинация][]RK{}

	for k, v := range комбинации_map_slice {
		достижимые_2[k] = append(достижимые_2[k], RK{
			In:  k,
			Out: v,
		})
	}
	for _, rk := range порождающие {
		достижимые_2[rk.In] = append(достижимые_2[rk.In], rk)
	}

	// получение всех стартовых
	startTerms := получениеСтартовых(fsa, cfg)

	новые_порождающие := map[Комбинация][]RK{}
	// step 3: удаление недостижимых

	новый_стартовый := Комбинация{
		A: startTerms[0].A,
	}
	// и объединение
	for _, стартовыйНетерминал := range startTerms {
		следующие_порождающие := получениеДостижимых(достижимые_2, стартовыйНетерминал)

		новые_порождающие[новый_стартовый] = append(новые_порождающие[новый_стартовый], RK{
			In: новый_стартовый,
			Out: []Комбинация{
				стартовыйНетерминал,
			}},
		)

		for k, v := range следующие_порождающие {
			новые_порождающие[k] = append(новые_порождающие[k], v...)
		}
	}

	fmt.Println("---")
	fmt.Println("После удаление недостижимых")
	for k := range новые_порождающие {
		print(новые_порождающие[k])
	}

	без_дубликатов := удалениеДубликатов(новые_порождающие)

	for k := range без_дубликатов {
		print(без_дубликатов[k])
	}

	return без_дубликатов
}

func получениеПростыхПравил(cfg CFG, fsa FSA) (map[Комбинация][]Комбинация, []Rule) {
	комбинации_map_slice := map[Комбинация][]Комбинация{}

	rules_2 := []Rule{}

	for _, rule := range cfg.Rules {
		// len(rule.terms) != 1 точно ли ?
		if len(rule.Nterms) != 0 || len(rule.terms) != 1 {
			rules_2 = append(rules_2, rule)
			continue
		}
		pp := fsa.PиAB[rule.terms[0]]

		for A, BB := range pp {
			for _, B := range BB {
				kk := Комбинация{
					qi: A,
					A:  rule.Nterm,
					qj: B,
				}
				комбинации_map_slice[kk] = append(комбинации_map_slice[kk], Комбинация{
					t: rule.terms[0],
				})
			}
		}
	}
	return комбинации_map_slice, rules_2
}

func получениеОстальныхПравил(fsa FSA, rules_2 []Rule) []RK {
	комбинации_2 := []RK{}

	// step 2: X -> Y,  X -> X
	// L:
	for _, rule := range rules_2 {
		// игнорирование переходов в себя
		// for _, nterm := range rule.Nterms {
		// 	if rule.Nterm == nterm {
		// 		rules_3 = append(rules_3, rule)
		// 		continue L
		// 	}
		// }
		// правила для всех возможных p, q, qi
		// <p, A, q> - > <p, A1, q1> <qn-1, An, q>
		// <X1, X, X2>
		lines := получитьЦепочки(fsa, len(rule.Nterms)+1)
		fmt.Println("---")
		fmt.Println(lines)
		fmt.Println("---")

		// построение цепочки
		for i := range lines {
			kkk := []Комбинация{}
			kk := Комбинация{
				qi: lines[i][0],
				A:  rule.Nterm,
				qj: lines[i][len(rule.Nterms)],
			}
			// -> <Y1, Y, Y2>
			for j, nterm := range rule.Nterms {
				kkk = append(kkk, Комбинация{
					qi: lines[i][j], //
					A:  nterm,
					qj: lines[i][j+1], //
				})
			}

			комбинации_2 = append(комбинации_2, RK{
				In:  kk,
				Out: kkk,
			})
		}
	}
	return комбинации_2
}

func получениеПорождающих(комбинации_map_slice map[Комбинация][]Комбинация, комбинации_2 []RK) []RK {
	m := -1

	достижимые_нетермы := make(map[Комбинация]struct{}, len(комбинации_map_slice))
	for k := range комбинации_map_slice {
		достижимые_нетермы[k] = struct{}{}
	}

	fmt.Println("---")
	fmt.Println(достижимые_нетермы)
	fmt.Println("---")
	print(комбинации_2)
	fmt.Println("---")

	недостижимые := комбинации_2
	достижимые := []RK{}

	for m != len(достижимые_нетермы) {
		m = len(достижимые_нетермы)
		достижимые_нетермы, недостижимые, достижимые = удалениеНепорождающих(достижимые_нетермы, недостижимые, достижимые)
	}

	fmt.Println("---")
	fmt.Println(достижимые_нетермы)
	fmt.Println("---")
	print(недостижимые)
	fmt.Println("---")
	print(достижимые)
	fmt.Println("---")

	return достижимые
}

func получитьЦепочки(fsa FSA, length int) [][]string {
	lines := [][]string{}

	for A := range fsa.ABиP {
		l, endLine := getLine(fsa, A, length-1, [][]string{{A}})

		if endLine {
			lines = append(lines, l...)
		}
	}

	return lines
}

func getLine(fsa FSA, first string, length int, lines [][]string) ([][]string, bool) {
	if length == 0 {
		return lines, true
	}

	BB := fsa.ABиP[first]

	if len(BB) == 0 {
		return lines, false
	}

	// if len(BB) > len(lines) {
	for i := 0; i < len(BB)-1; i++ {
		lines = append(lines, lines[0])
	}
	// }
	i := 0

	new_lines := [][]string{}

	for bb := range BB {
		lines[i] = append(lines[i], bb)
		l, endLine := getLine(fsa, bb, length-1, [][]string{lines[i]})

		if endLine {
			new_lines = append(new_lines, l...)
		}
		i++
	}

	return new_lines, true
}

// комбинация -> []комбинации
type RK struct {
	In  Комбинация
	Out []Комбинация
}

type Комбинация struct {
	qi string
	A  string
	qj string

	t string
}

func удалениеНепорождающих(достижимые_нетермы map[Комбинация]struct{}, недостижимые, достижимые []RK,
) (достижимые_нетермы_2 map[Комбинация]struct{}, недостижимые_2, достижимые_2 []RK) {
	достижимые_нетермы_2 = достижимые_нетермы
	достижимые_2 = достижимые

C:
	for _, комбинация := range недостижимые {
		for _, nterm := range комбинация.Out {
			if _, ok := достижимые_нетермы[nterm]; !ok {
				недостижимые_2 = append(недостижимые_2, комбинация)
				continue C
			}
		}

		достижимые_2 = append(достижимые_2, комбинация)
		достижимые_нетермы[комбинация.In] = struct{}{}
	}

	return достижимые_нетермы, недостижимые_2, достижимые_2
}

func получениеДостижимых(порождающие map[Комбинация][]RK, стартовыйНетерминал Комбинация,
) map[Комбинация][]RK {
	следующиеНетерминалы := []Комбинация{}
	новые_порождающие := map[Комбинация][]RK{}

	достижимые_3 := map[Комбинация][]RK{}
	for k, v := range порождающие {
		достижимые_3[k] = v
	}

	for _, rk := range достижимые_3[стартовыйНетерминал] {
		следующиеНетерминалы = append(следующиеНетерминалы, rk.Out...)
		новые_порождающие[rk.In] = append(новые_порождающие[rk.In], rk)
	}

	delete(достижимые_3, стартовыйНетерминал)

	for len(следующиеНетерминалы) != 0 {
		номые_нетерминалы := []Комбинация{}
		for _, next_rks := range следующиеНетерминалы {
			for _, rk := range достижимые_3[next_rks] {
				номые_нетерминалы = append(номые_нетерминалы, rk.Out...)
				новые_порождающие[rk.In] = append(новые_порождающие[rk.In], rk)
			}

			delete(достижимые_3, next_rks)
		}

		следующиеНетерминалы = номые_нетерминалы
	}

	return новые_порождающие
}

func получениеСтартовых(fsa FSA, cfg CFG) []Комбинация {
	startTerms := []Комбинация{}

	for _, state := range fsa.EndStates {
		startTerms = append(startTerms, Комбинация{
			qi: fsa.FirstState,
			A:  cfg.StarmNTerm,
			qj: state,
		})
	}

	return startTerms
}

func удалениеДубликатов(StoAN map[Комбинация][]RK) map[Комбинация][]RK {
	AtoS := map[Комбинация]map[Комбинация]struct{}{}

	for ком, мRK := range StoAN {
		for _, rk := range мRK {
			for _, out := range rk.Out {
				if _, ok := AtoS[out]; !ok {
					AtoS[out] = map[Комбинация]struct{}{}
				}
				AtoS[out][ком] = struct{}{}
			}
		}
	}

	мапа_для_дубликатов := make(map[string]Комбинация)

	m := -1

	for len(StoAN) != m {
		m = len(StoAN)

		for _, mRK := range StoAN {
			for _, rkStoAn := range mRK {
				new_terms := getStringRule(rkStoAn.Out)

				old_терм, ok := мапа_для_дубликатов[new_terms]
				if !ok || old_терм == rkStoAn.In {
					мапа_для_дубликатов[new_terms] = rkStoAn.In
					continue
				}

				rkCДубликатами := AtoS[rkStoAn.In]

				for Sдубликат := range rkCДубликатами {
					mapПолучившихсяДубликатов := map[string]int{}

					поДубликаты := StoAN[Sдубликат]

					for i, rk := range поДубликаты {
						s_out := getStringRule(rk.Out)

						mapПолучившихсяДубликатов[s_out] = i

						for k, outtttt := range rk.Out {
							if outtttt == rkStoAn.In {
								rk.Out[k] = old_терм
							}
						}

						s_new_out := getStringRule(rk.Out)

						number, ok := mapПолучившихсяДубликатов[s_new_out]

						if !ok || number == i {
							continue
						}

						StoAN[Sдубликат] = append(StoAN[Sдубликат][:i], StoAN[Sдубликат][i+1:]...)
					}
				}

				// delete(AtoS[rkStoAn.In], rkStoAn.In)

				if len(StoAN[rkStoAn.In]) == 1 {
					delete(StoAN, rkStoAn.In)
				}

			}
		}

		fmt.Println("---")
		for k := range StoAN {
			print(StoAN[k])
		}
		fmt.Println("---")
	}

	// S -> AN : A -> S -> {S -> AN}
	// ntern_в_правилах := map[Комбинация]map[Комбинация][]*RK{}

	// for _, v := range порождающие {
	// 	for _, vv := range v {
	// 		rule := &vv
	// 		for _, out := range vv.Out {
	// 			if _, ok := ntern_в_правилах[out]; !ok {
	// 				ntern_в_правилах[out] = map[Комбинация][]*RK{}
	// 			}

	// 			ntern_в_правилах[out][vv.In] = append(ntern_в_правилах[out][vv.In], rule)
	// 		}
	// 	}
	// }

	// m := -1

	// for m != len(ntern_в_правилах) {
	// 	m = len(ntern_в_правилах)

	// 	мапа_для_дубликатов := make(map[string]Комбинация)

	// 	// A -> S -> []{S -> AN}
	// 	for A, m_k_rules := range ntern_в_правилах {
	// 		// S -> []{S -> AN}
	// 		for _, rules := range m_k_rules {
	// 			// S -> AN
	// 			for _, rk := range rules {
	// 				terms := getStringRule(rk.Out)

	// 				// <qi, A, qj><qi, N, qj> -> S
	// 				старый_терм, ok := мапа_для_дубликатов[terms]

	// 				if !ok || старый_терм == rk.In {
	// 					мапа_для_дубликатов[terms] = rk.In

	// 					// for _, out := range rk.Out {
	// 					// 	if _, ok := next_rules[out]; !ok {
	// 					// 		next_rules[out] = map[Комбинация][]*RK{}
	// 					// 	}
	// 					// 	next_rules[out][rk.In] = append(next_rules[out][rk.In], rk)
	// 					// }
	// 					continue
	// 				}
	// 				// находим все дубликаты -> заменяем

	// 				// map[S] -> []{S -> AN}
	// 				rks := ntern_в_правилах[rk.In]
	// 				for S, rkssss := range rks {
	// 					// []{S -> AN}
	// 					for j := range rkssss {
	// 						for k, outtttt := range rkssss[j].Out {
	// 							if outtttt == rk.In {
	// 								rkssss[j].Out[k] = старый_терм
	// 							}
	// 						}
	// 					}
	// 					ntern_в_правилах[старый_терм][S] = append(ntern_в_правилах[старый_терм][S], rkssss...)
	// 				}
	// 				// rk.In - дубликат

	// 				delete(ntern_в_правилах, A)
	// 			}
	// 		}
	// 	}
	// }

	return StoAN
}

func getStringRule(out []Комбинация) string {
	s := ""
	for _, v := range out {
		s += "<" + v.qi + v.A + v.qj + v.t + ">"
	}

	return s
}
