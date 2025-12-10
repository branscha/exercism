package railfence

import (
    "strings"

)

func Encode(message string, rails int) string {
	b := make([]strings.Builder, rails)
    pos := 0
    d := -1
    
    // Snake walk
    for _, r := range message {
        b[pos].WriteRune(r)
        if pos == 0 || pos == rails-1 {
            d *= -1
        }
        pos +=d
    }
    // Collect rails
    bc := strings.Builder{}
    for _, bldr := range b {
        bc.WriteString(bldr.String())
    }
    return bc.String()
}

func Decode(message string, rails int) string {
    msg := []rune(message)
    msglen := len(msg)
    blocks := make([][]rune, rails)

    // Split into rails
	p := (2*rails) - 2 // Zig-zag phase
    l1 := (msglen + p - 1)/p  // Letter frequency on outer rails
    l2 := msglen/(p/2) // Letter frequency on internal rails

    pos := 0
    for r := 0; r < rails; r++ {
        if r == 0 || r == rails-1 {
            if pos+l1 < msglen {
                blocks[r] = msg[pos:pos+l1]
            } else {
                blocks[r] = msg[pos:]
            }
            pos += l1
        } else {
            if pos+l2 < msglen {
                blocks[r] = msg[pos:pos+l2]
            } else {
                blocks[r] = msg[pos:]
            }
            pos += l2
        } 
    }

    // Snake walk
    sb := strings.Builder{}
	row := 0
    d := -1
    for i := 0; i < msglen; i++ {
        sb.WriteRune(blocks[row][0])
        blocks[row] = blocks[row][1:]
        if row == 0 || row == rails-1{
             d *= -1
        }
        row += d
    }
    return sb.String()
}
