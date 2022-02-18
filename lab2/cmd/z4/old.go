package main

// star(ε)=ε, e.ε=e, ∅+e=e, ∅.e=∅
func рекурсия(dfa DFA) DFA {
	for A, B := range dfa.ABиP {
		dfa.ABиP[A] = ОбъединениеПереходов(B)
	}

	dfa = добавление_нач_и_кон_состояния(dfa)

	printlnE(dfa)

	// for len(dfa.ABиP) > 2 {
	// 	// если в вершину есть вход и выход из одной вершины, то сначала избавляемся от нее
	// 	// dfa = удалениеПетель(dfa)

	// 	dfa = удалениеПромежуточных(dfa)
	// }

	// for A, B := range dfa.ABиP {
	// 	dfa.ABиP[A] = ОбъединениеПереходов(B)
	// }

	// решение номер 2
	dfa = решение_2(dfa)

	for A, B := range dfa.ABиP {
		dfa.ABиP[A] = ОбъединениеПереходов(B)
	}

	return dfa
}

// старое решение
// // после ОбъединениеПереходов остается только один в pp (pp = все переходы между двумя вершинами)
// func удалениеПромежуточных(dfa DFA) DFA {
// 	for A, RR := range dfa.ABиP {
// 		printlnDFA(dfa)

// 		for V1, V2 := range dfa.ABиP {
// 			dfa.ABиP[V1] = ОбъединениеПереходов(V2)
// 		}

// 		for R, pp1 := range RR {
// 			if A == R || R == dfa.КонечноеСостояние {
// 				continue
// 			}

// 			i := 0
// 			var естьпереходыВR bool

// 			// fmt.Println("A: ", A)
// 			// fmt.Println("R:", R)
// 			// fmt.Println(dfa.ABиP[R])
// 			// берем вершины А и В, удаляя между ними переходы R
// 			for B, pp2 := range dfa.ABиP[R] {
// 				if B == R {
// 					continue
// 				}

// 				// если в среднее состояние ещё есть переходы из другой вершины
// 				// тогда мы его оставляем
// 				// if len(dfa.BAиP[R]) > 0 {
// 				// 	for a1 := range dfa.BAиP[R] {
// 				// 		if a1 != A {
// 				// 			естьпереходыВR = true
// 				// 		}
// 				// 	}
// 				// }

// 				newP := ""
// 				if ppR, ok := dfa.ABиP[R][R]; ok {
// 					newP = "(" + pp1[0] + "(" + ppR[0] + ")*" + pp2[0] + ")"
// 					delete(dfa.ABиP[R], R)
// 				} else {
// 					if pp1[0] == "ε" {
// 						newP = pp2[0]
// 					} else if pp2[0] == "ε" {
// 						newP = pp1[0]
// 					} else {
// 						newP = "(" + pp1[0] + pp2[0] + ")"
// 					}

// 				}

// 				if !естьпереходыВR {
// 					delete(dfa.ABиP[R], B)
// 				}

// 				dfa.ABиP[A][B] = append(dfa.ABиP[A][B], newP)

// 				i++
// 			}

// 			if i > 0 {
// 				delete(dfa.ABиP[A], R)

// 				if pp3, ok := dfa.ABиP[R]; (!ok || len(pp3) == 0) && !естьпереходыВR {
// 					delete(dfa.ABиP, R)
// 				}
// 			}
// 		}
// 	}

// 	return dfa
// }
