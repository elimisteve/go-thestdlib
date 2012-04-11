package main

import "log"

func main() {
    ints := []int{1, 2, 3, 4, 5, 6}
    otherInts := []int{11, 12, 13, 14, 15, 16}

    log.Println(ints)
    log.Println(otherInts)

    copied := copy(ints[:3], otherInts)
    log.Printf("Copied %d ints from otherInts to ints", copied)

    log.Println(ints)
    log.Println(otherInts)
}
