package main

/*
#cgo CFLAGS: -mmacosx-version-min=10.6 -D__MAC_OS_X_VERSION_MAX_ALLOWED=1060
#cgo LDFLAGS: -framework CoreFoundation -framework Security
#include "keychain.h"

extern struct keychain_items* list_keychain_items();

*/
import "C"

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"
	"unsafe"
)

type AlfredOutputItems struct {
	XMLName xml.Name           `xml:"items"`
	Items   []AlfredOutputItem `xml:"item"`
}

type AlfredOutputItem struct {
	XMLName      xml.Name `xml:"item"`
	UID          string   `xml:"uid,attr"`
	Autocomplete string   `xml:"autocomplete,attr"`
	Title        string   `xml:"title"`
	Arg          string   `xml:"arg"`
}

func main() {
	var pattern string
	if len(os.Args) > 1 {
		pattern = strings.ToUpper(os.Args[1])
	}

	items := filterStringsIgnoreCase(keychainItems(), pattern)
	alfredItems := toAlfredItems(items)
	rawXML := toXML(alfredItems)
	fmt.Println(rawXML)
}

func toXML(alfredOut AlfredOutputItems) string {
	xmlout, err := xml.MarshalIndent(alfredOut, "", "  ")
	if err != nil {
		panic(err)
	}
	xmlWithHeader := xml.Header + string(xmlout)
	return xmlWithHeader
}

func toAlfredItems(items []string) AlfredOutputItems {
	alfredOut := AlfredOutputItems{Items: make([]AlfredOutputItem, 0, 10)}
	for _, item := range items {
		alfredOut.Items = append(alfredOut.Items, AlfredOutputItem{
			UID:          "kc:" + item,
			Autocomplete: item,
			Title:        item,
			Arg:          item,
		})
	}
	return alfredOut
}

type match struct {
	s   string
	pos int
}

type matchList []match

func (lst matchList) Len() int {
	return len(lst)
}

func (lst matchList) Swap(i, j int) {
	lst[i], lst[j] = lst[j], lst[i]
}

func (lst matchList) Less(i, j int) bool {
	return lst[i].pos < lst[j].pos
}

func filterStringsIgnoreCase(strs []string, pattern string) []string {
	var matches = make(matchList, 0, len(strs))
	for _, s := range strs {
		matchPos := strings.Index(strings.ToUpper(s), pattern)
		if matchPos == -1 {
			continue
		}
		matches = append(matches, match{s, matchPos})
	}
	sort.Sort(matches)

	var result = make([]string, 0, len(strs))
	for _, m := range matches {
		result = append(result, m.s)
	}
	return result
}

func keychainItems() []string {
	var itemsPtr *C.struct_keychain_items = C.list_keychain_items()
	length := int(itemsPtr.Count)
	var items = make([]string, length)

	defer C.free(unsafe.Pointer(itemsPtr.Items))
	defer C.free(unsafe.Pointer(itemsPtr))

	cstrings := (*[1 << 30]*C.char)(unsafe.Pointer(itemsPtr.Items))[:length:length]

	for _, s := range cstrings {
		str := C.GoString(s)
		defer C.free(unsafe.Pointer(s))
		items = append(items, str)
	}
	return items
}

func CFStringToStr(cs C.CFStringRef) string {
	valcstr := C.CFStringGetCStringPtr(cs, C.kCFStringEncodingUTF8)
	return C.GoString(valcstr)
}
