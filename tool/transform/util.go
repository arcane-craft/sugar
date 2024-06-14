package main

import (
	"crypto/sha256"
	"encoding/base32"
	"go/types"
	"path"
	"regexp"
	"strings"
)

var mainTypeRe = regexp.MustCompile(`^\**([a-zA-Z0-9_\-./]+)`)
var anyTypeRe = regexp.MustCompile(`[a-zA-Z0-9_\-./]+`)

func GetPkgPathFromTypeStr(str string) string {
	matches := mainTypeRe.FindStringSubmatch(str)
	if len(matches) == 2 {
		typ := matches[1]
		ext := path.Ext(typ)
		if len(ext) > 0 {
			return strings.TrimSuffix(typ, ext)
		}
	}
	return ""
}

func GetPkgPathFromType(ty types.Type) string {
	return GetPkgPathFromTypeStr(ty.String())
}

func GetNameFromTypeStr(str string) string {
	matches := mainTypeRe.FindStringSubmatch(str)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

func GetNameFromType(ty types.Type) string {
	return GetNameFromTypeStr(ty.String())
}

func GetTypesFromTypeStr(str string) []string {
	return anyTypeRe.FindAllString(str, -1)
}

func GetTypesFromType(ty types.Type) []string {
	return GetTypesFromTypeStr(ty.String())
}

func ResetTypeStrPkgName(str string, imports map[string]string, currentPkg string) (ret string, adds map[string]string) {
	typs := GetTypesFromTypeStr(str)
	ret = str
	for _, typ := range typs {
		pkgPath := GetPkgPathFromTypeStr(typ)
		if pkgPath == currentPkg {
			if len(pkgPath) > 0 {
				ret = strings.ReplaceAll(ret, pkgPath+".", "")
			}
		} else if len(pkgPath) > 0 {
			pkgName, ok := imports[pkgPath]
			if !ok {
				pkgName = GetRandPkgName(path.Base(pkgPath), pkgPath)
				if adds == nil {
					adds = map[string]string{}
				}
				adds[pkgPath] = pkgName
			}
			if len(pkgName) > 0 {
				pkgName += "."
			}
			ret = strings.ReplaceAll(ret, pkgPath+".", pkgName)
		}
	}
	return
}

func ResetTypePkgName(ty types.Type, imports map[string]string, currentPkg string) (ret string, adds map[string]string) {
	ret, adds = ResetTypeStrPkgName(ty.String(), imports, currentPkg)
	return
}

func GetRandPkgName(prefix string, nameHint string) string {
	sum := sha256.Sum256([]byte(nameHint))
	return prefix + "_" + strings.TrimRight(base32.HexEncoding.EncodeToString(sum[13:19]), "=")
}

func GetRandVarName(nameHint string) string {
	sum := sha256.Sum256([]byte(nameHint))
	return "var" + strings.TrimRight(base32.HexEncoding.EncodeToString(sum[13:19]), "=")
}
