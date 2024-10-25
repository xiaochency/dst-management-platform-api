package home

func Success(message string, lang string) string {
	successZH := map[string]string{
		"rollbackSuccess": "回档成功",
		"restartSuccess":  "重启成功",
		"shutdownSuccess": "关闭成功",
		"startupSuccess":  "开启成功",
		"updating":        "正在更新中，请耐心等待",
		"announceSuccess": "宣告成功",
		"execSuccess":     "执行成功",
		"resetSuccess":    "重置成功",
	}
	successEN := map[string]string{
		"rollbackSuccess": "Rollback Success",
		"restartSuccess":  "Restart Success",
		"shutdownSuccess": "Shutdown Success",
		"startupSuccess":  "Startup Success",
		"updating":        "Updating, please wait patiently",
		"announceSuccess": "Announce Success",
		"execSuccess":     "Execute Success",
		"resetSuccess":    "Reset Success",
	}

	if lang == "zh" {
		return successZH[message]
	} else {
		return successEN[message]
	}
}
