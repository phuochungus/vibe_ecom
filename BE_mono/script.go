package main

import (
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	println(string(hash))

	// err := bcrypt.CompareHashAndPassword([]byte("$2a$10$dFRrJ38YZB7rxbwwPj1WVOGxvt8ZiHDCGUtGojIQ4LKJA0AgbTvcq"), []byte("123456"))
	// if err != nil {
	// 	panic(err)
	// }
	// println("Password is correct")
}
