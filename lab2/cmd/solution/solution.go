package solution

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type EquationsSystem struct {
	States     map[string]Equation
	Переменные map[string]struct{}

	Порядок []string
}

type Equation struct {
	Var   string
	Regex string
}

func GetData(configPath string) (EquationsSystem, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	data := strings.Split(strings.ReplaceAll(string(file), " ", ""), "\n")

	equationSystem := EquationsSystem{
		Переменные: map[string]struct{}{},
		States:     map[string]Equation{},
	}

	re := regexp.MustCompile(`\([a-z|ε]+`)
	re2 := regexp.MustCompile(`[a-z|ε]+\)`)
	re3 := regexp.MustCompile(`[^[a-z|ε|\)]]*[A-Z]`)
	re4 := regexp.MustCompile(`[A-Z]`)

	for _, v := range data {
		if v == "" {
			continue
		}

		parts := strings.SplitN(v, "=", 2)

		снова := parts[1]
		// проверка на уровень вложенности
		for strings.Contains(снова, "(") {
			убрано_все_до_скобки := strings.SplitN(снова, "(", 2)
			вторая_скобка := strings.Index(убрано_все_до_скобки[1], "(")
			плюс := strings.Index(убрано_все_до_скобки[1], "+")
			закрывающая_скобка := strings.Index(убрано_все_до_скобки[1], ")")

			if вторая_скобка != -1 && вторая_скобка < закрывающая_скобка &&
				вторая_скобка < плюс {
				fmt.Println("Много скобок")
			}

			убано_все_после_закрывающей := strings.SplitN(убрано_все_до_скобки[1], ")", 2)
			снова = убано_все_после_закрывающей[1]
		}

		// проверка на альтерантиву
		альтернатива := parts[1]
		for strings.Contains(альтернатива, "|") {
			разделено_альтернативой := strings.SplitN(альтернатива, "|", 2)
			первая_половина := разделено_альтернативой[0]
			вторая_половина := разделено_альтернативой[1]

			скобки_до := re.FindAllString(первая_половина, -1)

			скопки_после := re2.FindAllString(вторая_половина, -1)

			if скобки_до == nil || скопки_после == nil {
				fmt.Println("Нет скобок")
			}

			убано_все_после_закрывающей := strings.SplitN(альтернатива, ")", 2)
			альтернатива = убано_все_после_закрывающей[1]
		}

		// проверка на наличие коэффициентов
		переменные_без_коэффициентов := re3.FindAllString("="+parts[1], -1)
		if len(переменные_без_коэффициентов) > 0 {
			fmt.Println("Переменные без коэффициентов")
		}

		// получение всех переменных в уравнении
		переменные := re4.FindAllString(parts[1], -1)
		for i := range переменные {
			equationSystem.Переменные[переменные[i]] = struct{}{}
		}

		equationSystem.States[parts[0]] = Equation{
			Var:   parts[0],
			Regex: parts[1],
		}

		equationSystem.Порядок = append([]string{parts[0]}, equationSystem.Порядок...)
	}

	for переменная := range equationSystem.Переменные {
		if _, ok := equationSystem.States[переменная]; !ok {
			fmt.Println("Нет уравления для:", переменная)
			os.Exit(0)
		}
	}

	return equationSystem, nil
}

var re_var = regexp.MustCompile(`[A-Z]`)
var re_var_and_letter = regexp.MustCompile(`(\([|a-z|ε]+\)|[a-z|ε]+)[A-Z]?`)
var re_letter = regexp.MustCompile(`(\([|a-z|ε]+\)|[a-z|ε]+)`)

func Pешение(equationSystem EquationsSystem) (ответ Equation) {
	// убираем X = aX
	for _, переменная := range equationSystem.Порядок {
		уравнение := equationSystem.States[переменная]

		if !strings.Contains(уравнение.Regex, переменная) {
			continue
		}

		all_expr := re_var_and_letter.FindAllString(уравнение.Regex, -1)
		for i := range all_expr {
			if strings.Contains(all_expr[i], переменная) {
				expr := re_letter.FindString(all_expr[i])

				if len(expr) > 1 {
					expr = "(" + expr + ")"
				}

				new_eq := expr + "*"
				if len(all_expr) > 1 {
					new_eq = new_eq + "("

					for j := range all_expr {
						if i == j {
							continue
						}

						new_eq = new_eq + all_expr[j] + "+"
					}
					new_eq = strings.TrimSuffix(new_eq, "+")
					new_eq = new_eq + ")"
				}

				уравнение.Regex = new_eq
				equationSystem.States[переменная] = уравнение

				break
			}
		}
	}

	for _, переменная := range equationSystem.Порядок {
		уравнение := equationSystem.States[переменная]

		for k, eq := range equationSystem.States {
			if переменная == k {
				continue
			}

			замена := уравнение.Regex
			if strings.Contains(уравнение.Regex, "+") {
				замена = "(" + уравнение.Regex + ")"
			}
			eq.Regex = strings.ReplaceAll(eq.Regex, переменная, замена)

			equationSystem.States[k] = eq
		}

		if len(equationSystem.States) > 1 {
			delete(equationSystem.States, переменная)
		} else {
			ответ = equationSystem.States[переменная]
		}
	}

	if strings.Contains(ответ.Regex, ответ.Var) {
		regex := ответ.Regex
		запас := ""
		if strings.Contains(ответ.Regex, "?") {
			sp := strings.Split(ответ.Regex, "?")
			запас = sp[1]
			regex = strings.TrimPrefix(strings.TrimSuffix(sp[0], ")"), "(")
		}

		all_expr := strings.Split(regex, "+")
		for i := range all_expr {
			if strings.Contains(all_expr[i], ответ.Var) {
				expr := strings.ReplaceAll(all_expr[i], ответ.Var, "")

				if len(expr) > 1 {
					expr = "(" + expr + ")"
				}

				new_eq := expr + "*"
				if len(all_expr) > 1 {
					new_eq = new_eq + "("

					for j := range all_expr {
						if i == j {
							continue
						}

						new_eq = new_eq + all_expr[j] + "+"
					}
					new_eq = strings.TrimSuffix(new_eq, "+")
					new_eq = new_eq + ")"
				}
				regex = new_eq

				break
			}
		}

		if запас != "" {
			regex = "(" + regex + ")?" + запас
		}

		ответ.Regex = regex
	}

	return ответ
}
