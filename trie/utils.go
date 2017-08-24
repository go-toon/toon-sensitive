package trie

import (
	"bufio"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var blackTrie *Trie
var whitePrefixTrie *Trie
var whiteSuffixTrie *Trie
var noiseWords *Noise

// InitAllTrie 初始化Trie库
func InitAllTrie() {
	BlackTrie()
	WhitePrefixTrie()
	WhiteSuffixTrie()
}

// BlackTrie 返回黑名单Trie树
func BlackTrie() *Trie {
	if blackTrie == nil {
		blackTrie = NewTrie()
		blackTrie.CheckWhiteList = true
		blackTrie.CheckNoise = true

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

func NoiseWords() *Noise {
	if noiseWords == nil {
		noiseWords = NewNoise()
		//initNoise(noiseWords, "./dicts/noise/default.txt")
		loadNoise(noiseWords, "./dicts/noise")

	}
	return noiseWords
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

// 读取脏词文件，初始化trie树
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

// 加载噪声词
func loadNoise(noise *Noise, path string) {
	var loadAllDictWalk filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		initNoise(noise, path)

		return nil
	}

	err := filepath.Walk(path, loadAllDictWalk)
	if err != nil {
		panic(err)
	}
}

// initNoise 初始化噪声词库
func initNoise(noise *Noise, path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("failed to open file %s %s", path, err.Error()))
	}

	defer f.Close()

	log.Printf("%s Load StopWords %s", time.Now().Local().Format("2006-01-02 15:04:05 -0700"), path)
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

		if word := string(line); word != "" {
			noise.Add(word)
		}
	}

	return
}
