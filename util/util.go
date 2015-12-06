package util

import(
  "math/rand"
)

var strValues = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numValues = []rune("0123456789")

// string random generation for all character
func RandString(size int) string {
    str := strValues
		str = append(str, numValues...)
    b := make([]rune, size)
    for i := range b {
        b[i] = str[rand.Intn(len(str))]
    }
    return string(b)
}

// string random generation for all number
func RandNum(size int) string {
    str := numValues
    b := make([]rune, size)
    for i := range b {
        b[i] = str[rand.Intn(len(str))]
    }
    return string(b)
}
