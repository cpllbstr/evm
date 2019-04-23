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
	verts     []int //кандидаты в группе
	cap       int   //максимальное число кандидатов в группе
	outdeltas []int
	filled    bool
}

func CreateGroups(n []int) []group {
	grs := []group{}
	for _, v := range n {
		grs = append(grs, group{cap: v, verts: []int{}, filled: false})
	}
	return grs
}

func (gp *group) Add(vert int) {
	if len(gp.verts) < gp.cap {
		gp.verts = append(gp.verts, vert)
	} else {
		panic("Adding more then capacity")
	}
	return
}

func (gp *group) CheckFilled() bool {
	if gp.cap == len(gp.verts) {
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

//Plot -  выводит текущее состояние матрицы смежности, вроде как красивенько даже
func (g graph) PlotSubm(subm []int) {
	fmt.Println()
	fmt.Print("  |")
	for _, v := range subm {
		fmt.Printf("% 3v", v)
	}
	fmt.Println()
	for range subm {
		fmt.Print("---")
	}
	fmt.Println("---")
	for _, i := range subm {
		fmt.Print(i, " |")
		for _, j := range subm {
			fmt.Printf("% 3v", g.distanses[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
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

// Neighbors - возвращает смежные c подмножеством вершины, ищет вершины по всему графу
func (g graph) Neighbors(verts, subm []int) []int {
	neigh := []int{}
	for _, vert := range verts { //вершины для которых ищем соседей
		for _, submvert := range subm { //в каком подмножестве
			if g.distanses[vert][submvert] > 0 {
				fl1, fl2 := CheckExist([]int{submvert}, neigh), CheckExist([]int{submvert}, verts)
				if fl1 && fl2 {
					neigh = append(neigh, submvert)
				}
			}
		}
	}
	return neigh
}

//Swap -  меняет местами узлы в группах
func Swap(vert1, vert2 int, grp []group) {
	g1i := -1
	v1i := -1
	g2i := -1
	v2i := -1
	for gi, g := range grp {
		for vi, v := range g.verts {
			if vert1 == v {
				g1i = gi
				v1i = vi
			}
			if vert2 == v {
				g2i = gi
				v2i = vi
			}
		}
	}
	buf := grp[g1i].verts[v1i]
	grp[g1i].verts[v1i] = grp[g2i].verts[v2i]
	grp[g2i].verts[v2i] = buf
}

func GropConfig(grp []group) []int {
	groupconf := []int{}
	for _, v := range grp {
		groupconf = append(groupconf, v.verts...)
	}
	return groupconf
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
	for _, i := range subm {
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
	min := 10000
	mini := 0
	maxi := 0
	for i, delta := range deltas {
		if delta < max {
			max = delta
			maxi = i
		}
		if delta > min {
			min = delta
			mini = i
		}
	}
	return ConfigWithout(cand, []int{cand[maxi]}), cand[mini]

}

func (g graph) CountNeigh(verts, subm []int) ([]int, int) {
	neigh := make([]int, len(verts))
	for i, vert := range verts { //вершины для которых ищем соседей
		for _, submvert := range subm { //в каком подмножестве
			if g.distanses[vert][submvert] > 0 {
				fl1 := CheckExist([]int{submvert}, verts)
				if fl1 {
					neigh[i]++
				}
			}
		}
	}
	min := 1000
	mini := 0
	for n, c := range neigh {
		if min > c {
			min = c
			mini = n
		}
	}
	return neigh, mini

}

func (gp *group) CalcDeltas(g graph) {
	gp.outdeltas = []int{}
	outer := ConfigWithout(g.config, gp.verts)
	for _, v := range gp.verts {
		gp.outdeltas = append(gp.outdeltas, g.Power(v, outer)-g.Power(v, gp.verts))
	}
}

func (g graph) VertsToSwap(grp []group) (int, int) {
	minlen := 1000
	mini := 0
	for i := range grp {
		if minlen > len(grp[i].verts) {
			mini = i
		}
	}
	max := 0
	maxx := -1
	maxy := -1
	for i := 0; i < len(grp); i++ {
		if i != mini {
			for x := range grp[mini].outdeltas {
				for y := range grp[i].outdeltas {
					Pi := grp[mini].outdeltas[x] + grp[i].outdeltas[y] - 2*g.Power(grp[mini].verts[x], []int{grp[i].verts[y]})
					if Pi > max {
						max = Pi
						maxx = grp[mini].verts[x]
						maxy = grp[i].verts[y]
					}
				}
			}

		}

	}
	fmt.Println(maxx, maxy)
	return maxx, maxy
}

func main() {
	divides := [][]int{{5, 5, 5, 5, 5, 5}, {6, 6, 6, 6, 6}, {7, 7, 5, 5, 6}}
	ngroups := divides[2]

	g := NewGraph(quest)
	//g.PlotSubm(g.config)
	grps := CreateGroups(ngroups)
	currconf := g.config

	//нахождение опорного графа
	for grupi := 0; grupi < len(grps); grupi++ {
		verts := g.FindSubMins(currconf)
		_, n := g.CountNeigh(verts, currconf)
		vertc := verts[n]
		cand := []int{vertc}
		for !grps[grupi].filled {
			if len(currconf) == grps[grupi].cap {
				grps[grupi].verts = currconf
				if grps[grupi].CheckFilled() {
					break
				}
			} else {
				for len(cand) < grps[grupi].cap {
					neigh := g.Neighbors(cand, currconf)
					_, bet := g.BetterCands(neigh, currconf)
					cand = append(cand, bet)
				}
				grps[grupi].verts = cand
				grps[grupi].CheckFilled()
			}
		}
		currconf = ConfigWithout(currconf, grps[grupi].verts)
	}

	//fmt.Println(grps)

	testmat := NewGraph(test2)
	tgrps := []group{
		group{cap: 4, verts: []int{0, 1, 2, 3}, outdeltas: []int{}, filled: true},
		group{cap: 4, verts: []int{4, 5, 6, 7}, outdeltas: []int{}, filled: true},
	}
	testmat.PlotSubm(testmat.config)
	for i := range tgrps {
		tgrps[i].CalcDeltas(testmat)
		fmt.Println(tgrps[i].outdeltas)
	}
	v1last := 0
	v2last := 0
	for v1, v2 := testmat.VertsToSwap(tgrps); v1 >= 0 && v2 >= 0; v1, v2 = testmat.VertsToSwap(tgrps) {
		fmt.Println("v1v2:", v1, v2)
		testmat.PlotSubm(GropConfig(tgrps))
		Swap(v1, v2, tgrps)
		for i := range tgrps {
			tgrps[i].CalcDeltas(testmat)
			fmt.Println(tgrps[i].outdeltas)
		}
		testmat.PlotSubm(GropConfig(tgrps))
		if v1last == v2 && v2last == v1 {
			break
		}
		v1last, v2last = v1, v2
	}

}