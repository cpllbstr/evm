package main

import "fmt"

/*
	В матрице хранятся значения расстояний на графе, она неизменна
	для перестановки меняются значения вектора config, в котором хранятся индексы строк матрицы
*/
type graph struct {
	distanses [][]int //расстояния в графе
	config    []int   //конфигурация графа
}

func VertexPower(v []int) int {
	p := 0
	for _, val := range v {
		p += val
	}
	return p
}

func NewGraph(matr [][]int) graph {
	var g graph
	for i := range matr {
		g.config = append(g.config, i)
	}
	g.distanses = matr
	return g
}

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
}

func (g graph) PrintMat() {
	for i, str := range g.distanses {
		fmt.Printf("%v|  %v -> %v\n", g.config[i], str, g.Power(i))
	}
}

//PowerGroup - считает глобальные степени вершин
func (g graph) PowerGroup(str, fin int) []int {
	var res []int
	for i := str; i <= fin; i++ {
		res = append(res, g.Power(i))
	}
	return res
}

//Power - степень вершины, для индексации текущей конфигурации
func (g graph) Power(index int) int {
	return VertexPower(g.distanses[g.config[index]])
}

//Swap -  меняет местами строки по ключам
func (g graph) Swap(in1, in2 int) {
	buf := g.config[in1]
	g.config[in1] = g.config[in2]
	g.config[in2] = buf
}

func (g graph) LocalPower(cand []int) {

}

func Delta(vert []int) {
	//Si := VertexPower(vert)

}

func main() {
	matr := test
	g := NewGraph(matr)
	fmt.Println(g.PowerGroup(0, 2))
}
