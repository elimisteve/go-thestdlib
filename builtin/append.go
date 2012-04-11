package main

import "log"

func main() {
    // Empty slice, with capacity of 10
    ints := make([]int, 0, 10)
    log.Println(ints)

    ints2 := append(ints, 1, 2, 3)

    log.Println(ints2)
    log.Printf("Slice was at %p, it's probably still at %p", ints, ints2)

    moreInts := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
    ints3 := append(ints2, moreInts...)

    log.Println(ints3)
    log.Printf("Slice was at %p, and it moved to %p", ints2, ints3)
}
