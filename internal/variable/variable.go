package variable

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var varPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)

func Replace(text string, envVars map[string]string, tempVars map[string]string) string {
	return varPattern.ReplaceAllStringFunc(text, func(match string) string {
		key := match[2 : len(match)-2]
		if tempVars != nil {
			if v, ok := tempVars[key]; ok {
				return v
			}
		}
		if envVars != nil {
			if v, ok := envVars[key]; ok {
				return v
			}
		}
		return match
	})
}

func ReplaceMap(m map[string]string, envVars map[string]string, tempVars map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = Replace(v, envVars, tempVars)
	}
	return result
}

func GenerateDynamic(funcName string, args string) string {
	switch funcName {
	case "timestamp":
		return strconv.FormatInt(time.Now().UnixMilli(), 10)
	case "unix":
		return strconv.FormatInt(time.Now().Unix(), 10)
	case "uuid":
		return generateUUID()
	case "random_string":
		length := 8
		if args != "" {
			if n, err := strconv.Atoi(args); err == nil && n > 0 {
				length = n
			}
		}
		return randomString(length)
	case "random_int":
		return strconv.Itoa(rand.Intn(10000))
	case "date":
		return time.Now().Format("2006-01-02")
	case "datetime":
		return time.Now().Format("2006-01-02 15:04:05")
	case "year":
		return strconv.Itoa(time.Now().Year())
	case "month":
		return strconv.Itoa(int(time.Now().Month()))
	case "day":
		return strconv.Itoa(time.Now().Day())
	case "hour":
		return strconv.Itoa(time.Now().Hour())
	case "minute":
		return strconv.Itoa(time.Now().Minute())
	case "second":
		return strconv.Itoa(time.Now().Second())
	default:
		return "{{" + funcName + "}}"
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return strings.Join([]string{
		hexBytes(b[0:4]),
		hexBytes(b[4:6]),
		hexBytes(b[6:8]),
		hexBytes(b[8:10]),
		hexBytes(b[10:]),
	}, "-")
}

func hexBytes(b []byte) string {
	const hex = "0123456789ABCDEF"
	var sb strings.Builder
	for _, by := range b {
		sb.WriteByte(hex[by>>4])
		sb.WriteByte(hex[by&0x0f])
	}
	return sb.String()
}

func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
