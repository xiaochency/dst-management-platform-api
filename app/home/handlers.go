package home

import (
	"dst-management-platform-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func handleRoomInfoGet(c *gin.Context) {

	type Data struct {
		RoomSettingBase utils.RoomSettingBase `json:"roomSettingBase"`
		SeasonInfo      metaInfo              `json:"seasonInfo"`
	}
	type Response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}
	filePath, err := findLatestMetaFile(utils.MetaPath)
	if err != nil {
		return
	}

	seasonInfo := getMetaInfo(filePath)
	config, _ := utils.ReadConfig()

	data := Data{
		RoomSettingBase: config.RoomSetting.Base,
		SeasonInfo:      seasonInfo,
	}

	response := Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}

	c.JSON(http.StatusOK, response)
}

func handleSystemInfoGet(c *gin.Context) {
	type Data struct {
		Cpu    float64 `json:"cpu"`
		Memory float64 `json:"memory"`
		Master int     `json:"master"`
		Caves  int     `json:"caves"`
	}
	type Response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}
	var response Response
	response.Code = 200
	response.Message = "success"
	response.Data.Cpu = utils.CpuUsage()
	response.Data.Memory = utils.MemoryUsage()
	response.Data.Master = getProcessStatus(utils.MasterScreenName)
	response.Data.Caves = getProcessStatus(utils.CavesScreenName)
	c.JSON(http.StatusOK, response)
}

func handleExecPost(c *gin.Context) {
	type ExecForm struct {
		Type string `json:"type"`
		Info int    `json:"info"`
	}
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}

	var execFrom ExecForm
	if err := c.ShouldBindJSON(&execFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch execFrom.Type {
	case "startup":
		masterStatus := getProcessStatus(utils.MasterScreenName)
		cavesStatus := getProcessStatus(utils.CavesScreenName)
		if masterStatus == 0 {
			_ = utils.BashCMD(utils.StartMasterCMD)
		}
		if cavesStatus == 0 {
			config, _ := utils.ReadConfig()
			if config.RoomSetting.Cave != "" {
				_ = utils.BashCMD(utils.StartCavesCMD)
			}
		}

		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("startupSuccess", langStr), "data": nil})

	case "rollback":
		cmd := "c_rollback(" + strconv.Itoa(execFrom.Info) + ")"
		err := utils.ScreenCMD(cmd, utils.MasterName)
		if err != nil {
			utils.RespondWithError(c, 511, langStr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("rollbackSuccess", langStr), "data": nil})

	case "shutdown":
		cmd := "c_shutdown()"
		_ = utils.ScreenCMD(cmd, utils.MasterName)
		_ = utils.ScreenCMD(cmd, utils.CavesName)
		time.Sleep(2 * time.Second)
		_ = utils.BashCMD(utils.StopMasterCMD)
		_ = utils.BashCMD(utils.StopCavesCMD)
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("shutdownSuccess", langStr), "data": nil})

	case "restart":
		config, _ := utils.ReadConfig()
		cmd := "c_shutdown()"
		_ = utils.ScreenCMD(cmd, utils.MasterName)
		if config.RoomSetting.Cave != "" {
			_ = utils.ScreenCMD(cmd, utils.CavesName)
		}

		time.Sleep(2 * time.Second)
		_ = utils.BashCMD(utils.StopMasterCMD)
		if config.RoomSetting.Cave != "" {
			_ = utils.BashCMD(utils.StopCavesCMD)
		}

		time.Sleep(1 * time.Second)
		_ = utils.BashCMD(utils.StartMasterCMD)
		if config.RoomSetting.Cave != "" {
			_ = utils.BashCMD(utils.StartCavesCMD)
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("restartSuccess", langStr), "data": nil})

	case "update":
		config, _ := utils.ReadConfig()
		cmd := "c_shutdown()"
		_ = utils.ScreenCMD(cmd, utils.MasterName)
		if config.RoomSetting.Cave != "" {
			_ = utils.ScreenCMD(cmd, utils.CavesName)
		}

		time.Sleep(2 * time.Second)
		_ = utils.BashCMD(utils.StopMasterCMD)
		if config.RoomSetting.Cave != "" {
			_ = utils.BashCMD(utils.StopCavesCMD)
		}

		go func() {
			_ = utils.BashCMD(utils.UpdateGameCMD)
			_ = utils.BashCMD(utils.StartMasterCMD)
			if config.RoomSetting.Cave != "" {
				_ = utils.BashCMD(utils.StartCavesCMD)
			}
		}()

		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("updating", langStr), "data": nil})

	case "reset":
		cmd := "c_regenerateworld()"
		_ = utils.ScreenCMD(cmd, utils.MasterName)

		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("resetSuccess", langStr), "data": nil})

	case "masterSwitch":
		if execFrom.Info == 0 {
			cmd := "c_shutdown()"
			_ = utils.ScreenCMD(cmd, utils.MasterName)
			time.Sleep(2 * time.Second)
			_ = utils.BashCMD(utils.StopMasterCMD)

			c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("shutdownSuccess", langStr), "data": nil})
		} else {
			//开启服务器
			_ = utils.BashCMD(utils.StartMasterCMD)
			c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("startupSuccess", langStr), "data": nil})
		}

	case "cavesSwitch":
		if execFrom.Info == 0 {
			cmd := "c_shutdown()"
			_ = utils.ScreenCMD(cmd, utils.CavesName)
			time.Sleep(2 * time.Second)
			_ = utils.BashCMD(utils.StopCavesCMD)

			c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("shutdownSuccess", langStr), "data": nil})
		} else {
			//开启服务器
			_ = utils.BashCMD(utils.StartCavesCMD)
			c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("startupSuccess", langStr), "data": nil})
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
}

func handleAnnouncementPost(c *gin.Context) {
	type AnnouncementForm struct {
		Message string `json:"message"`
	}
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}

	var announcementForm AnnouncementForm
	if err := c.ShouldBindJSON(&announcementForm); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := "c_announce('" + announcementForm.Message + "')"
	_ = utils.ScreenCMD(cmd, utils.MasterName)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("announceSuccess", langStr), "data": nil})
}

func handleConsolePost(c *gin.Context) {
	type ConsoleForm struct {
		CMD   string `json:"cmd"`
		World string `json:"world"`
	}
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var consoleForm ConsoleForm
	if err := c.ShouldBindJSON(&consoleForm); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd := consoleForm.CMD
	if consoleForm.World == "master" {
		_ = utils.ScreenCMD(cmd, utils.MasterName)
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("execSuccess", langStr), "data": nil})
		return
	}
	if consoleForm.World == "caves" {
		_ = utils.ScreenCMD(cmd, utils.CavesName)
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": Success("execSuccess", langStr), "data": nil})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
}
