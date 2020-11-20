package gofpdf

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type utf8fontfile struct {
	file         *utf8FontFile
	originalSize int64
}

var (
	utf8fontCache      = map[string]*utf8fontfile{}
	utf8fontCacheMutex = sync.RWMutex{}
)

func loadUTF8FontFromCache(fontPath string) (utf8f *utf8FontFile, originalSize int64, err error) {
	if fontPath == "" {
		err = errors.New("need fontPath")
		return
	}

	utf8fontCacheMutex.RLock()
	cache := utf8fontCache[fontPath]
	utf8fontCacheMutex.RUnlock()
	if cache != nil {
		utf8f = copyUTF8FontFile(cache.file)
		originalSize = cache.originalSize
		return
	}

	loaded, err := loadUTF8FontFromFile(fontPath)
	if err != nil {
		return
	}
	if loaded == nil {
		err = errors.New("can't load font from " + fontPath)
		return
	}
	utf8fontCacheMutex.Lock()
	utf8fontCache[fontPath] = loaded
	utf8fontCacheMutex.Unlock()

	utf8f = copyUTF8FontFile(loaded.file)
	originalSize = loaded.originalSize
	return
}

func loadUTF8FontFromFile(fontPath string) (utf8f *utf8fontfile, err error) {
	if fontPath == "" {
		err = errors.New("need fontPath")
		return
	}
	var ttfStat os.FileInfo
	ttfStat, err = os.Stat(fontPath)
	if err != nil {
		return
	}
	utf8f = &utf8fontfile{}
	utf8f.originalSize = ttfStat.Size()
	var utf8Bytes []byte
	utf8Bytes, err = ioutil.ReadFile(fontPath)
	if err != nil {
		return
	}
	reader := fileReader{readerPosition: 0, array: utf8Bytes}
	utf8File := newUTF8Font(&reader)
	err = utf8File.parseFile()
	if err != nil {
		return
	}
	utf8f.file = utf8File
	return
}

func copyUTF8FontFile(utf8f *utf8FontFile) (ret *utf8FontFile) {
	if utf8f == nil {
		return
	}

	utf8fCopy := *utf8f
	clone := &utf8fCopy

	if clone.fileReader != nil {
		fr := *clone.fileReader
		clone.fileReader = &fr
	}

	if clone.tableDescriptions != nil {
		tableDescriptions := map[string]*tableDescription{}
		for k, v := range clone.tableDescriptions {
			if v != nil {
				cvCopy := *v
				cv := &cvCopy
				cv.checksum = make([]int, len(v.checksum))
				copy(cv.checksum, v.checksum)
				tableDescriptions[k] = cv
			}
		}
		clone.tableDescriptions = tableDescriptions
	}

	if clone.outTablesData != nil {
		outTablesData := map[string][]byte{}
		for k, v := range clone.outTablesData {
			cv := make([]byte, len(v))
			copy(cv, v)
			outTablesData[k] = cv
		}
		clone.outTablesData = outTablesData
	}

	if clone.symbolPosition != nil {
		symbolPosition := make([]int, len(clone.symbolPosition))
		copy(symbolPosition, clone.symbolPosition)
		clone.symbolPosition = symbolPosition
	}

	if clone.charSymbolDictionary != nil {
		charSymbolDictionary := map[int]int{}
		for k, v := range clone.charSymbolDictionary {
			charSymbolDictionary[k] = v
		}
		clone.charSymbolDictionary = charSymbolDictionary
	}

	if clone.CharWidths != nil {
		charWidths := make([]int, len(clone.CharWidths))
		copy(charWidths, clone.CharWidths)
		clone.CharWidths = charWidths
	}

	if clone.symbolData != nil {
		symbolData := map[int]map[string][]int{}
		for k, v := range clone.symbolData {
			cv := map[string][]int{}
			for k2, v2 := range v {
				cv2 := make([]int, len(v2))
				copy(cv2, v2)
				cv[k2] = cv2
			}
			symbolData[k] = cv
		}
		clone.symbolData = symbolData
	}

	if clone.CodeSymbolDictionary != nil {
		codeSymbolDictionary := map[int]int{}
		for k, v := range clone.CodeSymbolDictionary {
			codeSymbolDictionary[k] = v
		}
		clone.CodeSymbolDictionary = codeSymbolDictionary
	}
	ret = clone
	return
}
