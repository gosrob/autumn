package pkginfo

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type PackageInfo struct {
	ImportPath string `json:"ImportPath"`
	Module     Module `json:"module"`
}
type Module struct {
	Path string `json:"path"`
	Dir  string `json:"dir"`
}

func GetFullPackage(fileDir string) PackageInfo {
	// 执行`go list`命令获取当前文件所在的完整包名称
	cmd := exec.Command("go", "list", "-e", "-json")
	if fileDir != "" {
		cmd.Dir = fileDir
	}
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("执行`go list`命令失败: %v", err)
	}

	// 解析`go list`命令输出的JSON结果
	var pkgInfo PackageInfo
	if err := json.Unmarshal(output, &pkgInfo); err != nil {
		log.Fatalf("解析`go list`命令输出失败: %v", err)
	}

	return pkgInfo
}

func GetPackageFromFilePath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		return parts[len(parts)-2]
	}
	return ""
}
