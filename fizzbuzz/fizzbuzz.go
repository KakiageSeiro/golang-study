package fizzbuzz

import "fmt"

// Run nが3で割り切れる場合はFizz、5で割り切れる場合はBuzz、15で割り切れる場合はFizzBuzz、そしてそれ以外の場合は、文字列に変換してそのまま返却する。
func Run(n uint32) string {
	switch {
	case n%15 == 0:
		return fmt.Sprint("FizzBuzz")
	case n%3 == 0:
		return fmt.Sprint("Fizz")
	case n%5 == 0:
		return fmt.Sprint("Buzz")
	default:
		return fmt.Sprintf("%d", n)
	}
}
