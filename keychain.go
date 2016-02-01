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

	items := filterKeychainItemsIgnoreCase(keychainItems(), pattern)
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

func toAlfredItems(items []keychainItem) AlfredOutputItems {
	alfredOut := AlfredOutputItems{Items: make([]AlfredOutputItem, 0, 10)}
	for _, item := range items {
		alfredOut.Items = append(alfredOut.Items, AlfredOutputItem{
			UID:          "kc:" + item.Label,
			Autocomplete: item.Label,
			Title:        fmt.Sprintf("%s (%s)", item.Label, item.Service),
			Arg:          item.Label,
		})
	}
	return alfredOut
}

type match struct {
	item keychainItem
	pos  int
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

func filterKeychainItemsIgnoreCase(items []keychainItem, pattern string) []keychainItem {
	var matches = make(matchList, 0, len(items))
	for _, item := range items {
		matchPos := strings.Index(strings.ToUpper(item.Label), pattern)
		if matchPos == -1 {
			continue
		}
		matches = append(matches, match{item, matchPos})
	}
	sort.Sort(matches)

	var result = make([]keychainItem, 0, len(items))
	for _, m := range matches {
		result = append(result, m.item)
	}
	return result
}

type keychainItem struct {
	Label   string
	Service string
}

func keychainItems() []keychainItem {
	itemsPtr := C.list_keychain_items()
	length := int(itemsPtr.Count)
	var items = make([]keychainItem, length)

	defer C.free(unsafe.Pointer(itemsPtr.Items))
	defer C.free(unsafe.Pointer(itemsPtr))

	citems := (*[1 << 30]*C.struct_keychain_item)(unsafe.Pointer(itemsPtr.Items))[:length:length]

	for _, citem := range citems {
		defer C.free(unsafe.Pointer(citem.Label))
		defer C.free(unsafe.Pointer(citem.Service))

		item := keychainItem{
			Label:   C.GoString(citem.Label),
			Service: C.GoString(citem.Service),
		}
		items = append(items, item)
	}
	return items
}
