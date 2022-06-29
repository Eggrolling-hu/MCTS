package main

import (
	"math"
	"mcts"
)

func main() {

	vSlice := []int64{
		23, 18, 22, 82, 61, 32, 72, 90, 56, 80, 25}
	wSlice := []int64{
		93, 81, 12, 98, 33, 73, 46, 32, 78, 40, 98}
	var capacity int64 = 300

	// vSlice := []int64{
	// 	90000, 89750, 10001, 89500, 10252, 89250, 10503, 89000, 10754, 88750,
	// 	11005, 88500, 11256, 88250, 11507, 88000, 11758, 87750, 12009, 87500,
	// 	12260, 87250, 12511, 87000, 12762, 86750, 13013, 86500, 13264, 86250}
	// wSlice := []int64{
	// 	90001, 89751, 10002, 89501, 10254, 89251, 10506, 89001, 10758, 88751,
	// 	11010, 88501, 11262, 88251, 11514, 88001, 11766, 87751, 12018, 87501,
	// 	12270, 87251, 12522, 87001, 12774, 86751, 13026, 86501, 13278, 86251}
	// var capacity int64 = 30000

	items := []Item{}
	for i := range vSlice {
		items = append(items, Item{wSlice[i], vSlice[i]})
	}

	s := KnapsackState{
		capacity:  capacity,
		weight:    0,
		value:     0,
		inside:    []Item{},
		leftItems: items,
	}

	mcts.UCTSearch(mcts.State(&s), int64(5e3), 100, 4e1*1/math.Sqrt(2))

	// t := mcts.State(&s)
	// fmt.Println(t.Log())
	// for i := 0; ; i++ {
	// 	fmt.Println("layers: ", i)
	// 	action := mcts.UCTSearch(mcts.State(&s), int64(5e5), 100, 1e4*1/math.Sqrt(2))
	// 	s.TakeAction(action)
	// 	if len(s.GetAvailableActions()) == 0 {
	// 		break
	// 	}
	// }
	// fmt.Println(s.Log())
	// fmt.Println(s.value, bestReward)

}
