package gotemplate

import (
	"bytes"
	"regexp"
	"strings"
	"sync"
	"text/template"
)

var flock = sync.Mutex{}

var fmap = template.FuncMap{
	"add":             add,
	"concat":          concat, // concat "a" "b" => "ab"
	"createMap":       createMap,
	"date":            dateFmt, // "2017-03-31 19:59:11" |  date "06.01.02" => "17.03.31"
	"dateFrom":        dateFmtLayout,
	"datetime":        datetime,
	"decimal":         decimalFmt, // 3.1415 decimal 6,2 => 3.14
	"div":             divide,
	"elseifthen":      notconditional, // elseifthen "a" "b" => b, elseifthen "" "b" => ""
	"empty":           empty,          // empty [] => "", ["bah"] => "bah"
	"escape":          escape,
	"explode":         explode,
	"filter":          filterPath,
	"fixlen":          fixlen,
	"fixlenr":         fixlenright,
	"float":           tofloat, // float "0123.234" => 123.234
	"formatUKDate":    formatUKDate,
	"ifthen":          conditional, // ifthen "a" "b" => a, ifthen "" "b" => b
	"in_array":        inArray,
	"int":             toint, // int "0123" => 123
	"isset":           isSet,
	"item":            item, // item "a:b" ":" 0 => a
	"json_decode":     jsonDecode,
	"json_encode":     jsonEncode,
	"json_escape":     jsonEscape,
	"json":            asJSON,
	"last":            last,
	"limit":           limit,
	"lower":           strings.ToLower,
	"mapto":           mapto, // mapto "a" "a:True|b:False" "|:" => True
	"match":           regexp.MatchString,
	"md5":             md5hash,
	"mkSlice":         mkSlice,
	"mul":             multiply,
	"nanotimestamp":   nanotimestamp,
	"replace":         replace,
	"reReplaceAll":    reReplaceAll,
	"sanitise":        sanitise,
	"sanitize":        sanitise,
	"seq":             seq,
	"setItem":         setItem,
	"sql":             sqlEscape,
	"sub":             subtract,
	"timeformat":      timeFormat,
	"timeformatminus": timeFormatMinus,
	"timestamp":       timestamp,
	"title":           strings.Title,
	"toAbs":           toAbs,
	"tojson":          jsonDecode, // backward compatibility
	"toLower":         strings.ToLower,
	"toUpper":         strings.ToUpper,
	"ukdate":          ukdate,
	"ukdatetime":      ukdatetime,
	"unique":          unique,
	"unixtimestamp":   unixtimestamp,
	"upper":           strings.ToUpper,
	"url_path":        urlPath, // SEO, Slugify
	"urldecode":       urldecode,
	"urlencode":       urlencode,
	"xml_array":       xmlArray,
	"xml_decode":      xmlDecode,
	"xml_encode":      xmlEncode,
}

// RegisterFunc registers a new template func to the template parser
func RegisterFunc(key string, templatefunc interface{}) {
	flock.Lock()
	defer flock.Unlock()
	fmap[key] = templatefunc
}

// GetFuncs will return all usable template funcs as string slice
func GetFuncs() []string {
	flock.Lock()
	defer flock.Unlock()
	keys := make([]string, 0, len(fmap))
	for k := range fmap {
		keys = append(keys, k)
	}
	return keys
}

// MustTemplate parses string as Go template, using data as scope
func MustTemplate(str string, data interface{}) string {
	ret, err := Template(str, data)
	if err != nil {
		panic(err)
	}
	return ret
}

// TemplateDelim parses string with custom delimiters as Go template, using data as scope
func TemplateDelim(str string, data interface{}, begin, end string) (string, error) {
	tmpl, err := template.New("test").Funcs(fmap).Delims(begin, end).Parse(str)
	if err == nil {
		var doc bytes.Buffer
		err = tmpl.Execute(&doc, data)
		if err != nil {
			return "", err
		}
		return strings.Replace(doc.String(), "<no value>", "", -1), nil
	}
	return "", err
}

// Template parses string as Go template, using data as scope
func Template(str string, data interface{}) (string, error) {
	tmpl, err := template.New("test").Funcs(fmap).Parse(str)
	if err == nil {
		var doc bytes.Buffer
		err = tmpl.Execute(&doc, data)
		if err != nil {
			return "", err
		}
		return strings.Replace(doc.String(), "<no value>", "", -1), nil
	}
	return "", err
}
