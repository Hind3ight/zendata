package gen

import (
	"fmt"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	shellUtils "github.com/easysoft/zendata/src/utils/shell"
	"github.com/easysoft/zendata/src/utils/vari"
	"path"
	"strconv"
	"strings"
)

const (
	outputDir = "output"
	bufFile   = "data.bin"
)

func GenerateProtobuf(protoFile string) (content, pth string) {
	outputDir := generateCls(protoFile)

	convertFile := generateConverter(outputDir)

	content, pth = generateBinData(convertFile)

	return
}

func generateBinData(convertFile string) (content, pth string) {
	dir := path.Dir(convertFile)

	phpExeFile := "php"
	if commonUtils.IsWin() { // use build-in php runtime
		phpExeFile = path.Join(vari.ZdPath, "runtime", "php7", "php.exe")
	}
	cmdStr := phpExeFile + " convert.php"
	out, _ := shellUtils.ExecInDir(cmdStr, dir)
	if vari.Verbose {
		logUtils.PrintTo(out)
	}

	pth = path.Join(dir, bufFile)

	return
}

func generateConverter(dir string) (pth string) {
	srcFile := path.Join(vari.ZdPath, "runtime", "protobuf", "convert.php")
	pth = path.Join(dir, "convert.php")

	content := fileUtils.ReadFile(srcFile)
	content = strings.ReplaceAll(content, "${cls_name}", vari.ProtoCls)

	fileUtils.WriteFile(pth, content)

	return
}

func generateCls(protoFile string) (ret string) {
	outputDir := path.Join(fileUtils.GetAbsoluteDir(protoFile), outputDir)
	fileUtils.RmFile(outputDir)
	fileUtils.MkDirIfNeeded(outputDir)

	platform := commonUtils.GetOs()
	execFile := "protoc"
	if commonUtils.IsWin() {
		platform += fmt.Sprintf("%d", strconv.IntSize)
		execFile += ".exe"
	}

	execFile = path.Join(vari.ZdPath, "runtime", "protobuf", "bin", platform, execFile)

	cmdStr := fmt.Sprintf("%s --php_out=%s %s", execFile, outputDir, protoFile)
	shellUtils.Exec(cmdStr)

	ret = outputDir

	return
}
