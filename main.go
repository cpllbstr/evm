package main

import (
	"fmt"
)

/*
	В матрице хранятся значения расстояний на графе, она неизменна
	для перестановки меняются значения вектора config, в котором хранятся индексы строк матрицы
*/
type graph struct {
	distanses [][]int //расстояния в графе
	config    []int   //конфигурация графа, внутри хранятся индексы изначального состояния графа, они же являются ключами
}

type group struct {
	cand   []int //кандидаты в группе
	cap    int   //максимальное число кандидатов в группе
	filled bool
}

func CreateGroups(n []int) []group {
	grs := []group{}
	for _, v := range n {
		grs = append(grs, group{cap: v, cand: []int{}, filled: false})
	}
	return grs
}

func (gp *group) Add(vert int) {
	if len(gp.cand) < gp.cap {
		gp.cand = append(gp.cand, vert)
	} else {
		panic("Adding more then capacity")
	}
	return
}

func (gp *group) CheckFilled() bool {
	if gp.cap == len(gp.cand) {
		gp.filled = true
	}
	return gp.filled
}

func NewGraph(matr [][]int) graph {
	var g graph
	for i := range matr {
		g.config = append(g.config, i)
	}
	g.distanses = matr
	return g
}

//summArr - суммирует числа в массиве
func summArr(v []int) int {
	p := 0
	for _, val := range v {
		p += val
	}
	return p
}

//Plot -  выводит текущее состояние матрицы смежности, вроде как красивенько даже
func (g graph) Plot() {
	fmt.Println()
	fmt.Print("  |")
	for _, v := range g.config {
		fmt.Printf("%3v", v)
	}
	fmt.Println()
	for range g.config {
		fmt.Print("---")
	}
	fmt.Println("---")
	for _, i := range g.config {
		fmt.Print(i, " |")
		for _, j := range g.config {
			fmt.Printf("%3v", g.distanses[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

/*
func (g graph) PrintMat() {
	for i, strt := range g.distanses {
		fmt.Printf("%v|  %v -> %v\n", g.config[i], strt, g.Power(i))
	}
}
*/

//PowerGroup - считает степени вершин
/*
	strt - индекс начального элемента диагональной подматрицы mat[strt][strt]
	finl - индекс конечного элемента
*/
func (g graph) PowerSubMatrix(subm []int) []int {
	var res []int
	for i := range subm {
		res = append(res, g.Power(i, subm))
	}
	return res
}

//Power - степень вершины в подматрице
func (g graph) Power(vert int, subm []int) int {
	pow := 0
	for _, v := range subm {
		pow += g.distanses[vert][v]
	}
	return pow
}

//CheckExist -  проверяет наличие в одной марице узлов из другой
func CheckExist(verts, conf []int) bool {
	for _, c := range conf {
		for _, v := range verts {
			if c == v {
				return false
			}
		}
	}
	return true
}

func ConfigWithout(conf, verts []int) []int {
	confw := []int{}
	fl := true
	for _, c := range conf {
		for _, v := range verts {
			if c == v {
				fl = false
				break
			}
		}
		if fl {
			confw = append(confw, c)
		}
		fl = true
	}
	return confw
}

// Neighbors - возвращает смежные c подмножеством вершины
func (g graph) Neighbors(verts []int) []int {
	neigh := []int{}
	for _, vert := range verts {
		for vertex, distanse := range g.distanses[vert] {
			if distanse > 0 {
				fl1, fl2 := CheckExist([]int{vertex}, neigh), CheckExist([]int{vertex}, verts)
				if fl1 && fl2 {
					neigh = append(neigh, vertex)
				}
			}
		}
	}
	return neigh
}

//Swap -  изменяет конфигурацию графа, меняет значения в матрице конфигурации
func (g graph) Swap(in1, in2 int) {
	buf := g.config[in1]
	g.config[in1] = g.config[in2]
	g.config[in2] = buf
}

//FindSubMin - ищет вершины с минимальной степенью в подмножестве. Возвращает номера вершин
func (g graph) FindSubMins(subm []int) []int {
	min := 1000
	minind := []int{}
	//поиск вершины с минимальной степенью
	for _, i := range subm {
		pow := g.Power(i, subm)
		if min > pow {
			min = pow
		}
	}
	for i := range g.distanses {
		pow := g.Power(i, subm)
		if min == pow {
			minind = append(minind, i)
		}
	}
	vert := []int{}
	for _, v := range minind {
		vert = append(vert, v)
	}
	return vert
}

func (g graph) BetterCands(cand, subm []int) ([]int, int) {
	deltas := []int{}
	for _, c := range cand {
		deltas = append(deltas, g.Power(c, subm)-g.Power(c, cand))
	}
	max := -10000
	maxi := 0
	for i, delta := range deltas {
		if delta < max {
			max = delta
			maxi = i
		}
	}
	fmt.Println("B:", cand, subm, cand[maxi])
	return ConfigWithout(cand, []int{cand[maxi]}), cand[maxi]

}

func main() {
	//ngroups := []int{5, 5, 5, 5, 5, 5}
	//ngroups := []int{6, 6, 6, 6, 6}
	//ngroups := []int{7, 7, 5, 5, 6}
	ngroups := []int{4, 4}
	g := NewGraph(test2)
	g.Plot()
	fmt.Println(g.FindSubMins([]int{0, 1, 2, 3, 4, 5, 6}))
	grps := CreateGroups(ngroups)
	currconf := g.config
	for grupi := 0; grupi < len(grps); {
		verts := g.FindSubMins(currconf)
		fmt.Println("verts:", verts, currconf)
		vertc := verts[0]
		neigh := g.Neighbors([]int{vertc})
		if len(neigh) == grps[grupi].cap {
			for _, n := range neigh {
				grps[grupi].Add(n)
			}
			grps[grupi].Add(vertc)
		}
		if len(neigh) > grps[grupi].cap {
			cand := append(neigh, vertc)
			for len(cand) != grps[grupi].cap {
				cand, _ = g.BetterCands(cand, currconf)
			}
			grps[grupi].cand = cand
		}
		/*if len(neigh) < grps[grupi].cap {
			fmt.Println(neigh)
			cand := append(neigh, vertc)
			for len(cand) != grps[grupi].cap {
				cand = g.BetterCands(cand, currconf)
			}
			grps[grupi].cand = cand
		}*/
		fmt.Println("Cand:", grps[grupi].cand)
		if grps[grupi].CheckFilled() {
			currconf = ConfigWithout(currconf, grps[grupi].cand)
			grupi++
		}
	}
	fmt.Println(grps)

}
