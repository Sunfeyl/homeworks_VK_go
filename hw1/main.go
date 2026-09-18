package main

import "fmt"

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		diff := target - nums[i]
		val, exists := m[diff]
		if exists {
			return []int{val, i}
		}
		m[nums[i]] = i
	}
	return []int{}
}
func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	ans := twoSum(nums, target)
	fmt.Println("Result:", ans)
}
