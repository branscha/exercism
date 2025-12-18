package change

import (
    "errors"
    "slices"
)

var changeError error = errors.New("Change error")

func Change(coins []int, target int) ([]int, error) {
    // Backtracking solution to find minimal coin combination
    // Backtrack position = nr used of each coin
    // Backtrack steps = vary the number of each coin used

    combo := make([]int, len(coins))
    comboCoins := 0
    comboSolutions := 0

    var try func([]int, []int, int, int)
    try = func(coins []int, state []int, left int, nc int) {
        // No coins left or already worse than best
        if len(coins) <= 0 || (nc >= comboCoins && comboCoins > 0) {
            return
        }

        p := len(coins) -1
        c := coins[p]
        for n := left/c; n>=0; n-- {
            
            state[p] = n
            left -= c * n
            nc += n
            if left == 0 {
                // solution, record it
                if comboCoins == 0 || nc < comboCoins {
                    comboCoins = nc
                    combo = slices.Clone(state)
                    comboSolutions += 1
                }
            } else {
                try(coins[:p], state, left, nc)
            }

            // Clean up 
        	state[p] = 0
            nc -= n
            left += c * n
            
        } // for candidates
        
    } // try

    try(coins, make([]int, len(coins)), target, 0)

    if comboSolutions > 0 {
        result := []int{}
        for i := 0; i < len(combo); i ++ {
            for j := 0; j < combo[i]; j++ {
                result = append(result, coins[i])
            }
        }
        return result, nil
        
    } else {
        return nil, changeError
    }

    
}