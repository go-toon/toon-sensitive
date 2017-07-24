package trie

import (
	"path/filepath"
	"os"
	"fmt"
	"github.com/labstack/gommon/log"
	"time"
	"bufio"
	"io"
	"strings"
)

var blackTrie *Trie
var whitePrefixTrie *Trie
var whiteSuffixTrie *Trie

// InitAllTrie 初始化Trie库
func InitAllTrie()  {
	BlackTrie()
	WhitePrefixTrie()
	WhiteSuffixTrie()
}

// BlackTrie 返回黑名单Trie树
func BlackTrie() *Trie {
	if blackTrie == nil {
		blackTrie = NewTrie()
		blackTrie.CheckWhiteList = true

		loadDict(blackTrie, "add", "./dicts/black/default")
		loadDict(blackTrie, "del", "./dicts/black/exclude")

	}
	return blackTrie
}

// WhitePrefixTrie 返回白名单前缀Trie树
func WhitePrefixTrie() *Trie {
	if whitePrefixTrie == nil {
		whitePrefixTrie = NewTrie()

		loadDict(whitePrefixTrie, "add", "./dicts/white/prefix")
	}

	return whitePrefixTrie
}

// WhiteSuffixTrie 返回白名单后缀Trie树
func WhiteSuffixTrie() *Trie {
	if whiteSuffixTrie == nil {
		whiteSuffixTrie = NewTrie()

		loadDict(whiteSuffixTrie, "add", "./dicts/white/suffix")
	}
	return whiteSuffixTrie
}

func ClearWhitePrefixTrie() {
	whitePrefixTrie = NewTrie()
}

func ClearWhiteSuffixTrie() {
	whiteSuffixTrie = NewTrie()
}

func loadDict(trieHandle *Trie, op, path string) {
	var loadAllDictWalk filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		initTrie(trieHandle, op, path)

		return nil
	}

	err := filepath.Walk(path, loadAllDictWalk)
	if err != nil {
		panic(err)
	}
}

func initTrie(trieHandle *Trie, op, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("failed to open file %s %s", path, err.Error()))
	}

	defer f.Close()

	log.Printf("%s Load dict: %s", time.Now().Local().Format("2006-01-02 15:04:05 -0700"), path)

	buf := bufio.NewReader(f)
	for {
		line, isPrefix, e := buf.ReadLine()
		if e != nil {
			if e != io.EOF {
				err = e
			}
			break
		}

		if isPrefix {
			continue
		}

		if word := strings.TrimSpace(string(line)); word != "" {
			tmp := strings.Split(word, " ")
			s := strings.Trim(tmp[0], " ")

			if "add" == op {
				trieHandle.Add(s)
			} else if "del" == op {
				trieHandle.Del(s)
			}
		}
	}

	return
}
