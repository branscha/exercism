package transpose

import (
    "strings"
)

func Transpose(input []string) []string {
    // Get the dimensions of the input
	rows := len(input)
    columns := 0
    for _, row := range input {
        if len(row) > columns {
            columns = len(row)
        }
    }

    // Allocate transposed slice of string builders
    t := make([]strings.Builder, columns)
    for i := 0; i < rows; i++ {
        for j := 0; j < columns; j++ {
            if j < len(input[i]) {
                // Add "forgotten" spaces lazily
                // We recognize missing gaps if the builder is smaller
                // than it should be for this row.
                for t[j].Len() < i {
                    t[j].WriteRune(' ')
                }
                t[j].WriteRune(rune(input[i][j]))
            } 
        }
    }

    // Convert the builders to strings
    result := []string{}
    for _, tr := range(t) {
        result = append(result, tr.String())
    }
    return result
}
