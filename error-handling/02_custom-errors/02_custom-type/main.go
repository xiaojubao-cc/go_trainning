package main

import "fmt"

type Customer struct {
	lat, long string
	err       error
}

func (n *Customer) Error() string {
	return fmt.Sprintf("a norgate math error occured: %v %v %v", n.lat, n.long, n.err)
}

func Sqrt(f float64) (float64, error) {
	if f < 0 {
		nme := fmt.Errorf("norgate math redux: square root of negative number: %v", f)
		return 0, &Customer{"3.14", "4.56", nme}
	}
	return 42, nil
}
func main() {
	_, err := Sqrt(-10)
	if err != nil {
		fmt.Println(err)
	}
}
