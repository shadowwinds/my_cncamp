package main

import "fmt"

func main() {
	arr := []string{"I", "am", "stupid", "and", "weak"}
	replaceMap := map[string]string{"stupid": "smart", "weak": "strong"}
	result := replaceString(arr, replaceMap)
	fmt.Println(result)
}

func replaceString(arr []string, replaceMap map[string]string) []string {
	for i, v := range arr {
		if replaceMap[v] != "" {
			arr[i] = replaceMap[v]
		}
	}
	return arr
}
