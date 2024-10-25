package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type RoomSettingBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GameMode    string `json:"gameMode"`
	PVP         bool   `json:"pvp"`
	PlayerNum   int    `json:"playerNum"`
	BackDays    int    `json:"backDays"`
	Vote        bool   `json:"vote"`
	Password    string `json:"password"`
	Token       string `json:"token"`
}

type RoomSetting struct {
	Base   RoomSettingBase `json:"base"`
	Ground string          `json:"ground"`
	Cave   string          `json:"cave"`
	Mod    string          `json:"mod"`
}

type AutoUpdate struct {
	Enable bool   `json:"enable"`
	Time   string `json:"time"`
}

type AutoAnnounce struct {
	Name      string `json:"name"`
	Enable    bool   `json:"enable"`
	Content   string `json:"content"`
	Frequency int    `json:"frequency"`
}

type AutoBackup struct {
	Enable bool   `json:"enable"`
	Time   string `json:"time"`
}

type Config struct {
	Username     string         `json:"username"`
	Nickname     string         `json:"nickname"`
	Password     string         `json:"password"`
	JwtSecret    string         `json:"jwtSecret"`
	RoomSetting  RoomSetting    `json:"roomSetting"`
	AutoUpdate   AutoUpdate     `json:"autoUpdate"`
	AutoAnnounce []AutoAnnounce `json:"autoAnnounce"`
	AutoBackup   AutoBackup     `json:"autoBackup"`
}

type OSInfo struct {
	Architecture    string
	OS              string
	CPUModel        string
	CPUCores        int32
	MemorySize      uint64
	Platform        string
	PlatformVersion string
	Uptime          uint64
}

func GenerateJWT(username string, jwtSecret []byte, expiration int) (string, error) {
	// 定义一个自定义的声明结构

	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiration) * time.Hour).Unix(), // 过期时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string, jwtSecret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func Sha512(input string) string {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	hashed := hasher.Sum(nil)
	return hex.EncodeToString(hashed)
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Base64Decode(input string) string {
	decodedData, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		fmt.Println("解码失败:", err)
		return ""
	}
	return string(decodedData)
}

func CreateConfig() {
	_, err := os.Stat("DstMP.sdb")
	if !os.IsNotExist(err) {
		return
	}
	var config Config
	config.Username = "admin"
	config.Password = "ba3253876aed6bc22d4a6ff53d8406c6ad864195ed144ab5c87621b6c233b548baeae6956df346ec8c17f5ea10f35ee3cbc514797ed7ddd3145464e2a0bab413"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 26
	randomString := make([]byte, length)
	for i := range randomString {
		// 从字符集中随机选择一个字符
		randomString[i] = charset[r.Intn(len(charset))]
	}
	config.JwtSecret = string(randomString)
	WriteConfig(config)
}

func ReadConfig() (Config, error) {
	content, _ := os.ReadFile("DstMP.sdb")
	//jsonData := Base64Decode(string(content))
	jsonData := string(content)
	var config Config
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		return Config{}, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return config, nil
}

func WriteConfig(config Config) {
	data, err := json.MarshalIndent(config, "", "    ") // 格式化输出
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	file, err := os.OpenFile("DstMP.sdb", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // 在函数结束时关闭文件
	// 写入 JSON 数据到文件
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func MWlang() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Request.Header.Get("X-I18n-Lang")
		c.Set("lang", lang)
	}
}

func MWtoken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")
		config, _ := ReadConfig()
		tokenSecret := config.JwtSecret
		_, err := ValidateJWT(token, []byte(tokenSecret))
		if err != nil {
			lang := c.Request.Header.Get("X-I18n-Lang")
			RespondWithError(c, 420, lang)
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetOSInfo() (*OSInfo, error) {
	architecture := runtime.GOARCH

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpuModel := cpuInfo[0].ModelName
	cpuCore := cpuInfo[0].Cores

	// 获取内存信息
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	memorySize := virtualMemory.Total

	// 获取主机信息
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	platformVersion := hostInfo.PlatformVersion
	platform := hostInfo.Platform
	uptime := hostInfo.Uptime
	osName := hostInfo.OS
	// 返回系统信息
	return &OSInfo{
		Architecture:    architecture,
		OS:              osName,
		CPUModel:        cpuModel,
		CPUCores:        cpuCore,
		MemorySize:      memorySize,
		Platform:        platform,
		Uptime:          uptime,
		PlatformVersion: platformVersion,
	}, nil
}

func TruncAndWriteFile(fileName string, fileContent string) {
	fileContentByte := []byte(fileContent)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("打开或创建文件时出错:", err)
		return
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 写入新数据
	_, err = file.Write(fileContentByte)
	if err != nil {
		fmt.Println("写入数据时出错:", err)
		return
	}
}

func DeleteDir(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("删除目录失败:", err)
		return
	}
}

func CpuUsage() float64 {
	// 获取 CPU 使用率
	percent, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("Error getting CPU percent: ", err)
		return 0
	}
	return percent[0]
}
func MemoryUsage() float64 {
	// 获取内存信息
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error getting virtual memory info: ", err)
		return 0
	}
	return vmStat.UsedPercent
}

func ScreenCMD(cmd string, world string) error {
	var totalCMD string
	if world == MasterName {
		totalCMD = "screen -S \"" + MasterScreenName + "\" -p 0 -X stuff \"" + cmd + "\\n\""
	}
	if world == CavesName {
		totalCMD = "screen -S \"" + MasterScreenName + "\" -p 0 -X stuff \"" + cmd + "\\n\""
	}

	cmdExec := exec.Command("/bin/bash", "-c", totalCMD)
	err := cmdExec.Run()
	if err != nil {
		return err
	}
	return nil
}

func BashCMD(cmd string) error {
	cmdExec := exec.Command("/bin/bash", "-c", cmd)
	err := cmdExec.Run()
	if err != nil {
		return err
	}
	return nil
}
