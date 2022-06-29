package main

import (
	"mcts"
	"strconv"
)

var bestReward float64 = 0.0

type Item struct {
	weight int64
	value  int64
}

type Pickup struct {
	item Item
}

func (p Pickup) Log() string {
	return "w:" + strconv.Itoa(int(
		p.item.weight)) + " v:" + strconv.Itoa(int(p.item.value))
}

type KnapsackState struct {
	capacity  int64
	weight    int64
	value     int64
	inside    []Item
	leftItems []Item
}

func (s *KnapsackState) GetAvailableActions() []mcts.Action {
	availableActions := make([]mcts.Action, 0, len(s.leftItems))
	for i := 0; i < len(s.leftItems); i++ {
		if s.weight+s.leftItems[i].weight <= s.capacity {
			availableActions = append(availableActions, Pickup{s.leftItems[i]})
		}
	}
	return availableActions
}

func (s *KnapsackState) TakeAction(action mcts.Action) {
	p := action.(Pickup)
	for i, v := range s.leftItems {
		if v == p.item {
			s.inside = append(s.inside, v)
			s.value += v.value
			s.weight += v.weight
			s.leftItems = append(s.leftItems[:i], s.leftItems[i+1:]...)
			break
		}
	}
}

func (s *KnapsackState) Evaluate() float64 {
	if s.value > int64(bestReward) {
		bestReward = float64(s.value)
	}
	return float64(s.value)
}

func (s *KnapsackState) IsForcedTerminated() bool {
	return false
}

func (s *KnapsackState) IsRivalRound() bool {
	return false
}

func (s *KnapsackState) Log() string {
	logInside := "Inside: "
	for _, item := range s.inside {
		logInside += "["
		logInside += Pickup{item}.Log()
		logInside += "],"
	}
	logLeft := "\nLeft: "
	for _, item := range s.leftItems {
		logLeft += "["
		logLeft += Pickup{item}.Log()
		logLeft += "],"
	}
	return logInside + logLeft
}

func (s *KnapsackState) Deepcopy() mcts.State {
	inside := make([]Item, len(s.inside))
	copy(inside, s.inside)
	leftItems := make([]Item, len(s.leftItems))
	copy(leftItems, s.leftItems)
	newKnapsackState := KnapsackState{
		capacity:  s.capacity,
		weight:    s.weight,
		value:     s.value,
		inside:    inside,
		leftItems: leftItems,
	}
	return &newKnapsackState
}
