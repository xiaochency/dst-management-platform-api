package setting

import (
	"dst-management-platform-api/app/externalApi"
	"dst-management-platform-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func handleRoomSettingGet(c *gin.Context) {
	config, err := utils.ReadConfig()
	if err != nil {
		utils.Logger.Error("配置文件读取失败", "err", err)
		utils.RespondWithError(c, 500, "zh")
		return
	}
	type Response struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    utils.RoomSetting `json:"data"`
	}
	response := Response{
		Code:    200,
		Message: "success",
		Data:    config.RoomSetting,
	}
	c.JSON(http.StatusOK, response)
}

func handleRoomSettingSavePost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var roomSetting utils.RoomSetting
	if err := c.ShouldBindJSON(&roomSetting); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config, err := utils.ReadConfig()
	if err != nil {
		utils.Logger.Error("配置文件读取失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}
	config.RoomSetting = roomSetting
	err = utils.WriteConfig(config)
	if err != nil {
		utils.Logger.Error("配置文件写入失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}

	err = saveSetting(config)
	if err != nil {
		utils.Logger.Error("房间配置保存失败", "err", err)
	}
	err = DstModsSetup()
	if err != nil {
		utils.Logger.Error("mod配置保存失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("saveFail", langStr), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("saveSuccess", langStr), "data": nil})
}

func handleRoomSettingSaveAndRestartPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var roomSetting utils.RoomSetting
	if err := c.ShouldBindJSON(&roomSetting); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config, err := utils.ReadConfig()
	if err != nil {
		utils.Logger.Error("配置文件读取失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}
	config.RoomSetting = roomSetting
	err = utils.WriteConfig(config)
	if err != nil {
		utils.Logger.Error("配置文件写入失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}

	err = saveSetting(config)
	if err != nil {
		utils.Logger.Error("房间配置保存失败", "err", err)
	}
	err = DstModsSetup()
	if err != nil {
		utils.Logger.Error("mod配置保存失败", "err", err)
	}
	restartWorld(c, config, langStr)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("restartSuccess", langStr), "data": nil})
}

func handleRoomSettingSaveAndGeneratePost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var roomSetting utils.RoomSetting
	if err := c.ShouldBindJSON(&roomSetting); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config, err := utils.ReadConfig()
	if err != nil {
		utils.Logger.Error("配置文件读取失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}
	config.RoomSetting = roomSetting
	err = utils.WriteConfig(config)
	if err != nil {
		utils.Logger.Error("配置文件写入失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}

	err = saveSetting(config)
	if err != nil {
		utils.Logger.Error("房间配置保存失败", "err", err)
	}
	err = DstModsSetup()
	if err != nil {
		utils.Logger.Error("mod配置保存失败", "err", err)
	}
	generateWorld(c, config, langStr)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("generateSuccess", langStr), "data": nil})

}

func handlePlayerListGet(c *gin.Context) {
	type PlayerList struct {
		Players   []utils.Players `json:"players"`
		AdminList []string        `json:"adminList"`
		BlockList []string        `json:"blockList"`
		WhiteList []string        `json:"whiteList"`
	}

	//config, err := utils.ReadConfig()
	//if err != nil {
	//	utils.Logger.Error("配置文件读取失败", "err", err)
	//	utils.RespondWithError(c, 500, "zh")
	//	return
	//}
	adminList := getList(utils.AdminListPath)
	blockList := getList(utils.BlockListPath)
	whiteList := getList(utils.WhiteListPath)

	var playList PlayerList
	//playList.Players = config.Players
	playList.Players = utils.STATISTICS[len(utils.STATISTICS)-1].Players
	playList.AdminList = adminList
	playList.BlockList = blockList
	playList.WhiteList = whiteList

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": playList})
}

func handleAdminAddPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var uidFrom UIDForm
	if err := c.ShouldBindJSON(&uidFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := addList(uidFrom.UID, utils.AdminListPath)
	if err != nil {
		utils.Logger.Error("添加管理员失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("addAdminFail", langStr), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("addAdmin", langStr), "data": nil})
}

func handleBlockAddPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var uidFrom UIDForm
	if err := c.ShouldBindJSON(&uidFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := addList(uidFrom.UID, utils.BlockListPath)
	if err != nil {
		utils.Logger.Error("添加黑名单失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("addBlockFail", langStr), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("addBlock", langStr), "data": nil})
}

func handleWhiteAddPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var uidFrom UIDForm
	if err := c.ShouldBindJSON(&uidFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := addList(uidFrom.UID, utils.WhiteListPath)
	if err != nil {
		utils.Logger.Error("添加白名单失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("addWhiteFail", langStr), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("addWhite", langStr), "data": nil})
}

func handleAdminDeletePost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var uidFrom UIDForm
	if err := c.ShouldBindJSON(&uidFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := deleteList(uidFrom.UID, utils.AdminListPath)
	if err != nil {
		utils.Logger.Error("删除管理员失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("deleteAdminFail", langStr), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("deleteAdmin", langStr), "data": nil})
}

func handleBlockDeletePost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var uidFrom UIDForm
	if err := c.ShouldBindJSON(&uidFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := deleteList(uidFrom.UID, utils.BlockListPath)
	if err != nil {
		utils.Logger.Error("删除黑名单失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("deleteBlockFail", langStr), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("deleteBlock", langStr), "data": nil})
}

func handleWhiteDeletePost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var uidFrom UIDForm
	if err := c.ShouldBindJSON(&uidFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := deleteList(uidFrom.UID, utils.WhiteListPath)
	if err != nil {
		utils.Logger.Error("删除白名单失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("deleteWhiteFail", langStr), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("deleteWhite", langStr), "data": nil})
}

func handleKick(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	var uidFrom UIDForm
	if err := c.ShouldBindJSON(&uidFrom); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errMaster := kick(uidFrom.UID, utils.MasterName)
	errCaves := kick(uidFrom.UID, utils.CavesName)

	if errMaster != nil && errCaves != nil {
		utils.Logger.Error("踢出玩家失败", "errMaster", errMaster, "errCaves", errCaves)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("kickFail", langStr), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("kickSuccess", langStr), "data": nil})
}

func handleImportFileUploadPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("uploadFail", langStr), "data": nil})
		return
	}
	//保存文件
	savePath := utils.ImportFileUploadPath + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.Logger.Error("文件保存失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("uploadFail", langStr), "data": nil})
		return
	}
	//检查导入文件是否合法
	result, err := checkZipFile(file.Filename)
	if err != nil {
		utils.Logger.Error("检查导入文件失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("wrongUploadFile", langStr), "data": nil})
		return
	}
	if !result {
		utils.Logger.Error("导入文件校验失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("wrongUploadFile", langStr), "data": nil})
		return
	}
	//关闭服务器
	_ = utils.StopGame()
	//备份服务器
	err = utils.BackupGame()
	if err != nil {
		utils.Logger.Warn("游戏备份失败", "err", err)
	}
	//删除旧服务器文件
	err = utils.BashCMD("rm -rf " + utils.ServerPath + "*")
	if err != nil {
		utils.Logger.Error("删除旧服务器文件失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("deleteOldServerFail", langStr), "data": nil})
		return
	}
	//创建新服务器文件
	err = utils.BashCMD("mv " + utils.ImportFileUnzipPath + "* " + utils.ServerPath)
	if err != nil {
		utils.Logger.Error("创建新服务器文件失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("createNewServerFail", langStr), "data": nil})
		return
	}
	//写入数据库
	err = WriteDatabase()
	if err != nil {
		utils.Logger.Error("导入文件写入数据库失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("writeToDBFail", langStr), "data": nil})
		return
	}
	//清理上传的文件
	clearUpZipFile()
	// 写入dedicated_server_mods_setup.lua
	err = DstModsSetup()
	if err != nil {
		utils.Logger.Error("mod配置保存失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("saveFail", langStr), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("uploadSuccess", langStr), "data": nil})
}

func handleModSettingFormatGet(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}

	luaScript, _ := utils.GetFileAllContent(utils.MasterModPath)

	modInfo, err := externalApi.GetModsInfo(luaScript, langStr)
	if err != nil {
		utils.RespondWithError(c, 500, langStr)
		return
	}

	var responseData []utils.ModFormattedData
	for _, i := range utils.ModOverridesToStruct(luaScript) {
		item := utils.ModFormattedData{
			ID: i.ID,
			Name: func() string {
				for _, j := range modInfo {
					if i.ID == j.ID {
						return j.Name
					}
				}
				return ""
			}(),
			Enable:               i.Enabled,
			ConfigurationOptions: i.ConfigurationOptions,
			PreviewUrl: func() string {
				for _, j := range modInfo {
					if i.ID == j.ID {
						return j.PreviewUrl
					}
				}
				return ""
			}(),
		}
		responseData = append(responseData, item)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": responseData})
}

func handleModConfigOptionsGet(c *gin.Context) {
	type ModConfigurationsForm struct {
		ID int `form:"id" json:"id"`
	}
	var modConfigurationsForm ModConfigurationsForm
	if err := c.ShouldBindQuery(&modConfigurationsForm); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	type ModConfig struct {
		ID            int                         `json:"id"`
		ConfigOptions []utils.ConfigurationOption `json:"configOptions"`
	}
	var modConfig ModConfig
	modID := modConfigurationsForm.ID
	modInfoLuaFile := utils.ModUgcPath + "/" + strconv.Itoa(modID) + "/modinfo.lua"
	isUgcMod, err := utils.FileDirectoryExists(modInfoLuaFile)
	if err != nil {
		utils.RespondWithError(c, 500, langStr)
		return
	}

	if !isUgcMod {
		modInfoLuaFile = utils.ModNoUgcPath + "/workshop-" + strconv.Itoa(modID) + "/modinfo.lua"
		exist, err := utils.FileDirectoryExists(modInfoLuaFile)
		if err != nil {
			utils.RespondWithError(c, 500, langStr)
			return
		}
		if !exist {
			c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("needDownload", langStr), "data": nil})
			return
		}
	}

	luaScript, _ := utils.GetFileAllContent(modInfoLuaFile)
	modConfig.ID = modID
	modConfig.ConfigOptions = utils.GetModConfigOptions(luaScript, langStr)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": modConfig})
}

func test(c *gin.Context) {
	type ModFormattedDataForm struct {
		ModFormattedData []utils.ModFormattedData `json:"modFormattedData"`
	}
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}

	config, err := utils.ReadConfig()
	if err != nil {
		utils.Logger.Error("配置文件读取失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}

	var modFormattedDataForm ModFormattedDataForm
	if err := c.ShouldBindJSON(&modFormattedDataForm); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	luaString := utils.ParseToLua(modFormattedDataForm.ModFormattedData)

	config.RoomSetting.Mod = luaString
	// Master/modoverrides.lua
	err = utils.TruncAndWriteFile(utils.MasterModPath, config.RoomSetting.Mod)
	if err != nil {
		utils.Logger.Error("MOD配置文件写入失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}
	if config.RoomSetting.Cave != "" {
		//Caves/modoverrides.lua
		err = utils.TruncAndWriteFile(utils.CavesModPath, config.RoomSetting.Mod)
		if err != nil {
			utils.Logger.Error("MOD配置文件写入失败", "err", err)
			utils.RespondWithError(c, 500, langStr)
			return
		}
	}

	err = utils.WriteConfig(config)
	if err != nil {
		utils.Logger.Error("配置文件写入失败", "err", err)
		utils.RespondWithError(c, 500, langStr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": nil})
}

func handleModDownloadPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}
	type ModDownloadForm struct {
		ID      int    `json:"id"`
		FileURL string `json:"file_url"`
	}
	var modDownloadForm ModDownloadForm
	if err := c.ShouldBindJSON(&modDownloadForm); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func() {
		if modDownloadForm.FileURL == "" {
			cmd := utils.GenerateModDownloadCMD(modDownloadForm.ID)
			err := utils.BashCMD(cmd)
			if err != nil {
				utils.Logger.Error("MOD下载失败", "err", err)
			}
		} else {
			err := externalApi.DownloadMod(modDownloadForm.FileURL, modDownloadForm.ID)
			if err != nil {
				utils.Logger.Error("MOD下载失败", "err", err)
			}
		}
	}()

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("downloading", langStr), "data": nil})
}

func handleSyncModPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}

	err := utils.SyncMods()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("syncModFail", langStr), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("syncModSuccess", langStr), "data": nil})
}

func handleDeleteDownloadedModPost(c *gin.Context) {
	lang, _ := c.Get("lang")
	langStr := "zh" // 默认语言
	if strLang, ok := lang.(string); ok {
		langStr = strLang
	}

	type DeleteForm struct {
		ISUGC bool `json:"isUgc"`
		ID    int  `json:"id"`
	}

	var deleteForm DeleteForm
	if err := c.ShouldBindJSON(&deleteForm); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := utils.DeleteDownloadedMod(deleteForm.ISUGC, deleteForm.ID)
	if err != nil {
		utils.Logger.Error("删除已下载的MOD失败", "err", err)
		c.JSON(http.StatusOK, gin.H{"code": 201, "message": response("deleteModFail", langStr), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": response("deleteModSuccess", langStr), "data": nil})
}
