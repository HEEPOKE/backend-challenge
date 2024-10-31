package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// โครงสร้างข้อมูลสำหรับการตอบกลับ JSON ที่เก็บจำนวนคำที่เกี่ยวข้องกับเนื้อ
type BeefCountResponse struct {
	Beef map[string]int `json:"beef"`
}

// ฟังก์ชันหาผลรวมของเส้นทางที่น้อยที่สุดในสามเหลี่ยม
func minPathSum(triangle [][]int) int {
	if len(triangle) == 0 {
		return 0
	}
	// สร้าง slice เพื่อเก็บค่าผลรวมเส้นทางที่น้อยที่สุด
	dp := make([]int, len(triangle))
	copy(dp, triangle[len(triangle)-1])
	// วนลูปจากแถวที่สองจากล่างสุดไปบนสุด
	for i := len(triangle) - 2; i >= 0; i-- {
		for j := 0; j < len(triangle[i]); j++ {
			dp[j] = triangle[i][j] + int(math.Min(float64(dp[j]), float64(dp[j+1])))
		}
	}
	return dp[0]
}

// ฟังก์ชันลบตัวเลขจากทางซ้ายที่มีค่ามากกว่าเลขถัดไป
func removeLeft(numbers *[]int) {
	for i := 0; i < len(*numbers)-1; i++ {
		if (*numbers)[i] > (*numbers)[i+1] {
			*numbers = append((*numbers)[:i], (*numbers)[i+1:]...)
			break
		}
	}
}

// ฟังก์ชันลบตัวเลขจากทางขวาที่มีค่ามากกว่าเลขทางซ้าย
func removeRight(numbers *[]int) {
	for i := len(*numbers) - 1; i > 0; i-- {
		if (*numbers)[i] > (*numbers)[i-1] {
			*numbers = append((*numbers)[:i], (*numbers)[i+1:]...)
			break
		}
	}
}

// ฟังก์ชันแปลง slice ของตัวเลขเป็นสตริง
func convertToString(numbers []int) string {
	var sb strings.Builder
	for _, num := range numbers {
		sb.WriteString(strconv.Itoa(num))
	}
	return sb.String()
}

// ฟังก์ชันหาค่าผลลัพธ์ที่น้อยที่สุดจากโค้ดที่กำหนด ('L' และ 'R')
func findMinimumSum(code string) string {
	numbers := []int{4, 1, 2, 3, 5}
	for _, c := range code {
		switch c {
		case 'L':
			removeLeft(&numbers)
		case 'R':
			removeRight(&numbers)
		}
	}
	return convertToString(numbers)
}

// ฟังก์ชันนับจำนวนคำที่เกี่ยวข้องกับเนื้อในข้อความ
func countBeefWords(text string) map[string]int {
	words := strings.Fields(strings.ToLower(text))
	beefCounts := make(map[string]int)
	// รายการคำที่เกี่ยวข้องกับเนื้อที่ต้องการนับ
	beefWords := []string{"fatback", "t-bone", "pastrami", "pork", "meatloaf", "jowl", "enim", "bresaola"}
	for _, word := range words {
		word = strings.Trim(word, ".,")
		for _, beefWord := range beefWords {
			if word == beefWord {
				beefCounts[word]++
			}
		}
	}
	return beefCounts
}

// ฟังก์ชันจัดการคำขอ HTTP สำหรับการสรุปจำนวนคำที่เกี่ยวข้องกับเนื้อ
func beefSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// เรียกข้อมูลจาก API Bacon Ipsum
	resp, err := http.Get("https://baconipsum.com/api/?type=meat-and-filler&paras=99&format=text")
	if err != nil {
		http.Error(w, "Failed to fetch data from Bacon Ipsum API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// อ่านข้อมูลจากการตอบกลับของ API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// นับจำนวนคำที่เกี่ยวข้องกับเนื้อ
	beefCounts := countBeefWords(string(body))

	// สร้างโครงสร้างข้อมูลสำหรับการตอบกลับ
	response := BeefCountResponse{
		Beef: beefCounts,
	}

	// เขียนข้อมูลการตอบกลับในรูปแบบ JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func main() {
	// เปิดไฟล์ JSON ที่มีข้อมูลสามเหลี่ยม
	file, err := os.Open("files/hard.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// อ่านข้อมูลจากไฟล์
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// แปลงข้อมูล JSON เป็นโครงสร้างข้อมูลสามเหลี่ยม
	var triangle [][]int
	err = json.Unmarshal(byteValue, &triangle)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// กำหนดข้อมูลสามเหลี่ยมเพิ่มเติมสำหรับการทดสอบ
	triangleTwo := [][]int{
		{59},
		{73, 41},
		{52, 40, 53},
		{26, 53, 6, 34},
	}

	// คำนวณผลรวมของเส้นทางที่น้อยที่สุดจากข้อมูลสามเหลี่ยม
	result := minPathSum(triangle)
	resultTwo := minPathSum(triangleTwo)

	fmt.Println("The minimum path sum is:", result)
	fmt.Println("The minimum path Two sum is:", resultTwo)

	// ทดสอบฟังก์ชัน findMinimumSum ด้วยชุดอินพุตตัวอย่าง
	testInputs := []string{"LLRR", "RLL", "LRLR", "RRLR"}
	for _, input := range testInputs {
		fmt.Printf("input = %s output = %s\n", input, findMinimumSum(input))
	}

	// ตั้งค่า HTTP Handler สำหรับเส้นทาง /beef/summary
	http.HandleFunc("/beef/summary", beefSummaryHandler)

	// เริ่มต้นเซิร์ฟเวอร์ HTTP ที่พอร์ต 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
