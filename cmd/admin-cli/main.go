package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "hash":
		hashPassword()
	case "check":
		checkPassword()
	case "token":
		generateToken()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Admin System CLI 工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  admin-cli hash      交互式生成密码哈希")
	fmt.Println("  admin-cli check     验证密码是否匹配哈希")
	fmt.Println("  admin-cli token     生成测试 JWT Token")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  admin-cli hash")
	fmt.Println("  输入密码: admin123")
	fmt.Println("  输出: $2a$10$...")
}

func hashPassword() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("输入密码: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if len(password) < 6 {
		log.Fatal("密码长度至少为 6 个字符")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("生成哈希失败: %v", err)
	}

	fmt.Println("")
	fmt.Println("密码哈希:")
	fmt.Println(string(hash))
	fmt.Println("")
	fmt.Println("SQL 插入语句:")
	fmt.Printf("INSERT INTO users (username, password, email, status) VALUES ('admin', '%s', 'admin@example.com', 'active');\n", string(hash))
}

func checkPassword() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("输入密码: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("输入哈希: ")
	hash, _ := reader.ReadString('\n')
	hash = strings.TrimSpace(hash)

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("")
		fmt.Println("❌ 密码不匹配")
	} else {
		fmt.Println("")
		fmt.Println("✓ 密码匹配")
	}
}

func generateToken() {
	fmt.Println("JWT Token 生成需要 JWT Secret")
	fmt.Println("请使用 API 登录接口获取 Token:")
	fmt.Println("")
	fmt.Println("curl -X POST http://localhost:8080/api/v1/auth/login \\")
	fmt.Println("  -H 'Content-Type: application/json' \\")
	fmt.Println("  -d '{\"username\":\"admin\",\"password\":\"admin123\"}'")
}
