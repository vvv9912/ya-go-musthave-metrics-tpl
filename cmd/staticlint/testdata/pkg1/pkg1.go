package pkg1

import "fmt"

func mulfunc(i int) (int, error) {
	return i * 2, nil
}

func errCheckFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	a := 4
	fmt.Printf("a=%f\n", a)
	mulfunc(5)           // want "expression returns unchecked error"
	res, _ := mulfunc(5) // want "assignment with unchecked error"
	fmt.Println(res)     // want "expression returns unchecked error"
}
