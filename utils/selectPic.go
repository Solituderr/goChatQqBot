package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func GetRandPic() string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 指定图片文件所在的文件夹路径
	rootpath, err := os.Getwd()
	content, err := os.ReadFile(rootpath + "\\picpath.txt")
	if err != nil {
		fmt.Println(err)
		return "g"
	}
	file := string(content)
	folderPath := file

	// 获取文件夹中的所有文件
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "g"
	}

	// 从文件列表中筛选出所有图片文件
	var imageFiles []string
	for _, file := range files {
		if isImage(file.Name()) {
			imageFiles = append(imageFiles, file.Name())
		}
	}

	// 随机选择一个图片文件
	if len(imageFiles) > 0 {
		selectedFile := imageFiles[rand.Intn(len(imageFiles))]
		fmt.Println("Selected file:", selectedFile)
		return folderPath + "\\" + selectedFile
	} else {
		return "g"
	}
}

func isImage(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}
