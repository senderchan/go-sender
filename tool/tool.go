package tool

import (
	"strings"
)

func SplitArgs(txt string) (res []string) {
	tmp := ""
	q := false
	for _, s := range []rune(txt) {
		a := string(s)
		if a == " " {
			if !q && len(tmp) > 0 {
				res = append(res, tmp)
				tmp = ""
				continue
			} else if !q && len(tmp) == 0 {
				continue
			}
		} else if a == "\"" {
			if !q && len(tmp) == 0 {
				q = true
				continue
			} else if q {
				if strings.HasSuffix(tmp, "\\") {
					tmp = strings.TrimRight(tmp, "\\")
				} else {
					q = false
					continue
				}
			} else if !q {
				if strings.HasSuffix(tmp, "\\") {
					tmp = strings.TrimRight(tmp, "\\")
				} else {
					q = true
					continue
				}
			}
		}
		tmp += a
	}
	if len(tmp) > 0 {
		res = append(res, tmp)
	}
	for i := range res {
		res[i] = strings.ReplaceAll(res[i], `\"`, `"`)
	}
	return
}

func ParseSignaling(s string) (name string, arg []string) {
	name = s
	if strings.Count(s, " ") > 0 {
		name = s[:strings.Index(s, " ")]
		arg = SplitArgs(s[strings.Index(s, " ")+1:])
	}
	return
}
