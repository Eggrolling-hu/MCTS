package mcts

import (
	"fmt"

	"github.com/rs/zerolog"
)

func UCTSearch(state State, iterationLimit int64, simulationLimit int64, ucbConstant float64) Action {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	uctLogger := GetLogger(false, true)

	root := NewTreeNode(nil, state, nil, ucbConstant)
	root.visits++

	for i := 0; i < int(iterationLimit); i++ {
		// -> TreePolicy
		node := root

		// Select
		for len(node.untriedActions) == 0 && len(node.children) > 0 {
			node = node.Select()
		}

		// Expand
		if len(node.untriedActions) > 0 {
			node = node.Expand()
		}

		// -> DefaultPolicy
		reward := node.Simulate(simulationLimit)
		uctLogger.Debug().Msgf("\033[1;36m iter: %4d | reward: %.4f \033[0m", i, reward)

		// Backpropagate
		node.Backpropagate(reward)
	}

	bestChild := root.GetBestChild()
	uctLogger.Info().Msgf("\033[1;32m chosen action: %s \033[0m", bestChild.action.Log())

	var stack []*TreeNode
	stack = append(stack, root)

	for len(stack) > 0 {
		node := stack[0]
		relationship := fmt.Sprintf("node: %p | parent: %p | ", node, node.parent)
		uctLogger.Info().Msgf(relationship + node.Log())
		if root.children != nil {
			stack = append(stack, node.children...)
		}
		stack = append(stack[:0], stack[1:]...)
	}

	return bestChild.action

}
