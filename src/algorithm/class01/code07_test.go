package class01

import "testing"

func Test_07(t *testing.T) {

}

//异或运算：arr中，有两个数a和b出现了奇数次，其他数出现了偶数次，找出a和b
func findTwoOfOddTimesNum(arr []int) (a, b int) {
	//temp1 == a ^ b
	temp1 := 0
	for _, val := range arr {
		temp1 = temp1 ^ val
	}

	//temp2 == temp1二进制值的最右边一个1保留，其它位置都是0，假设该位置为x
	temp2 := temp1 & (-temp1)

	//arr被分为两组，一组x位置为1（假设a位于该组），一组x位置为0（b位于该组）
	//temp3 == a
	temp3 := 0
	for _, val := range arr {
		if (val & temp2) != 0 {
			temp3 = temp3 ^ val
		}
	}

	a = temp3
	b = temp3 ^ temp1
	return a, b
}
