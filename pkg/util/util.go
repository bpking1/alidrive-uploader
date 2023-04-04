package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"math"
	//"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

func GetSha1Code(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func GetProofCode(accessToken string, realpath string) (ProofCode, error) {
	var proofCode = ProofCode{
		Sha1:      "DA39A3EE5E6B4B0D3255BFEF95601890AFD80709",
		ProofCode: "",
	}
	stat, err := os.Stat(realpath)
	if err != nil {
		return proofCode, nil
	}
	filesize := stat.Size()
	if filesize == 0 {
		return proofCode, nil
	}
	file, err := os.Open(realpath)
	if err != nil {
		return proofCode, err
	}
	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return proofCode, err
	}
	proofCode.Sha1 = hex.EncodeToString(h.Sum(nil))
	escape := url.QueryEscape(accessToken)
	if err != nil {
		return proofCode, err
	}
	buffa := html.UnescapeString(escape)
	m := md5.New()
	m.Write([]byte(buffa))
	hashMd5 := hex.EncodeToString(m.Sum(nil))
	parseInt, err := strconv.ParseUint(hashMd5[0:16], 16, 64)
	if err != nil {
		return ProofCode{}, err
	}
	start := int64(parseInt % uint64(filesize))
	end := int64(math.Min(float64(start+8), float64(filesize)))

	_, err = file.Seek(start, 0)

	if err != nil {
		return ProofCode{}, err
	}
	var buffb = make([]byte, end-start)
	_, _ = file.Read(buffb)
	encoding := base64.StdEncoding.EncodeToString(buffb)
	proofCode.ProofCode = encoding
	return proofCode, nil
}
func GetFileContentType(out *os.File) string {
	// buffer := make([]byte, 512)
	// _, err := out.Seek(0, 0)
	// if err != nil {
	// 	return "plain/text"
	// }
	// _, err = out.Read(buffer)

	// defer func() { out.Seek(0, 0) }()
	// if err != nil {
	// 	return "plain/text"
	// }
	// contentType := http.DetectContentType(buffer)
	// return contentType
	return "application/octet-stream"
}

func GetAllFiles(path string, matchPattern string) ([]string, error) {
	var (
		files      []string
		matchRegex *regexp.Regexp = nil
		err        error
	)
	if matchPattern != "" {
		matchRegex, err = regexp.Compile(matchPattern)
		if err != nil {
			return files, err
		}
	}
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if matchRegex != nil && !matchRegex.MatchString(filepath.Base(p)) {
			return nil
		}
		files = append(files, filepath.ToSlash(p)[len(filepath.Dir(path)):])
		return nil
	})
	sort.Slice(files, func(i, j int) bool {
		return len(files[i]) < len(files[j])
	})
	return files, err
}

func TruncateText(str string, max int) string {
	length := len(str)
	if length-max > 3 {
		l := int(math.Min(float64(max), float64(length)))
		ru := []rune(str)
		str = string(ru[:l]) + "..."
	}
	return str
}

func FileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// FormatFileSize 字节的单位转换 保留两位小数
func FormatFileSize(fileSize float64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", fileSize/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", fileSize/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", fileSize/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", fileSize/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", fileSize/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", fileSize/float64(1024*1024*1024*1024*1024))
	}
}
