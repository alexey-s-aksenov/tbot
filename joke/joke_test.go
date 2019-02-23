package joke

import (
	"strconv"
	"testing"
)

func TestGetTheBest(t *testing.T) {
	testArray := make([]*quote, 0)
	for _, s := range []int{3, 5, 4, 2, 1} {
		testArray = append(testArray, &quote{int64(s), strconv.Itoa(s)})
	}
	actual := getTheBest(testArray)
	expected := "5"
	if actual != expected {
		t.Errorf("Wrong result. Expected: %v. Received: %v.\nTest array:\n%v", expected, actual, testArray)
	}
}
