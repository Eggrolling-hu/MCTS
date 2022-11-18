package mcts

import (
	"fmt"
	"math/rand"
	"sort"
)

// Action is the interface to change the states.
type Action interface {
	Log() string
}

// State is the interface to describe a game supports to satisfy the MCTS.
type State interface {
	Deepcopy() State
	GetAvailableActions() []Action // Return all the available actions
	TakeAction(action Action)      // Take action to change the state
	IsForcedTerminated() bool      // Return the state terminate for any special limitation
	Evaluate() float64             // Return the state evaluation
	IsRivalRound() bool            // Return the state min-max selection
	Log() string
}

// TreeNode is the node definition in the MCTS.
type TreeNode struct {
	parent         *TreeNode   // What node contains this node? Root node's parent is nil.
	action         Action      // What move lead to this node? Root node's action is nil.
	state          State       // What is the game state at this node?
	totalReward    float64     // What is the sum of all outcomes computed for this node and its children? From the point of view of a single player.
	visits         int64       // How many times has this node been studied? Used with totalValue to compute an average value for the node.
	untriedActions []Action    // What moves have not yet been explored from this state?
	children       []*TreeNode // The children of this node, can be many.
	ucbConstant    float64     // The UCB constant used in selection calculation.
	ucbValues      float64     // The computed score for this node used in selection, balanced between exploitation and exploration.
	IsRivalRound   bool        // For min-max tree to control the selection
}

func NewTreeNode(parent *TreeNode, state State, action Action, ucbConstant float64) *TreeNode {
	var node TreeNode = TreeNode{
		parent:         parent,
		action:         action,
		state:          state,
		totalReward:    0.0,                         // No outcome yet.
		visits:         0,                           // No visits yet.
		untriedActions: state.GetAvailableActions(), // Initially the node starts with every node unexplored.
		children:       nil,                         // No children yet.
		ucbConstant:    ucbConstant,                 // Whole tree uses same constant.
		ucbValues:      0.0,                         // No value yet.
		IsRivalRound:   state.IsRivalRound(),        // For min-max tree to control the selection
	}
	return &node
}

func (n *TreeNode) Select() *TreeNode {
	sort.Sort(byValues(n.children))
	return n.children[0]
}

func (n *TreeNode) Expand() *TreeNode {
	var i int = rand.Intn(len(n.untriedActions))
	var action Action = n.untriedActions[i]

	n.untriedActions = append(n.untriedActions[:i], n.untriedActions[i+1:]...)

	var newState State = n.state.Deepcopy()
	newState.TakeAction(action)

	var child *TreeNode = NewTreeNode(n, newState, action, n.ucbConstant)
	n.children = append(n.children, child)

	return child
}

func (n *TreeNode) Simulate(simulationLimit int64) float64 {

	var simulatedState State = n.state.Deepcopy()

	for i := 0; i < int(simulationLimit); i++ {
		var availableActions []Action = simulatedState.GetAvailableActions()
		if len(availableActions) == 0 {
			break
		}

		var i int = rand.Intn(len(availableActions))
		var action Action = availableActions[i]
		simulatedState.TakeAction(action)

		if simulatedState.IsForcedTerminated() {
			break
		}
	}

	return simulatedState.Evaluate()
}

func (n *TreeNode) Backpropagate(reward float64) {
	// Allow the root to call this on its parent with no ill effect.
	if n != nil {
		// Update this node's data.
		if n.IsRivalRound {
			n.totalReward -= reward
		} else {
			n.totalReward += reward
		}
		// n.totalReward += reward
		n.visits++
		// Recurse up the tree to the root.
		n.parent.Backpropagate(reward)

		var parentVisits int64 = 0
		if n.parent != nil {
			parentVisits = n.parent.visits
		}

		n.ucbValues = UCB(n.totalReward, n.ucbConstant, parentVisits, n.visits)
	}
}

func (n *TreeNode) GetBestChild() *TreeNode {
	sort.Sort(byVisits(n.children))
	// sort.Sort(byAverageReward(n.children))
	return n.children[0]
}

func (n *TreeNode) Log() string {
	var nodeLogging string
	nodeLogging += fmt.Sprintf(
		"visits: %6d | ucbValues: %9.2f | avgRewards: %9.2f | exploreValue: %9.2f",
		n.visits, n.ucbValues, n.totalReward/float64(n.visits),
		n.ucbValues-n.totalReward/float64(n.visits))
	return nodeLogging
}
