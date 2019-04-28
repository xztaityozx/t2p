package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildString(t *testing.T) {
	as := assert.New(t)

	// 引数有りのときのテストケース
	{
		type TestData struct {
			desc        string // テストの目的、説明
			args        []string
			exe         bool
			paletteSize int
			w, h        int
			box         []string
		}
		tds := []TestData{
			{desc: "引数を1つ指定する", args: []string{"hello"}, paletteSize: 128, w: 5, h: 1, box: []string{"hello"}},
			{desc: "引数を2つ指定する", args: []string{"hello", "world"}, paletteSize: 128, w: 11, h: 1, box: []string{"hello world"}},
			{desc: "引数に改行文字を含む", args: []string{"hello\nworld2"}, paletteSize: 128, w: 6, h: 2, box: []string{"hello", "world2"}},
			{desc: "シェルのコマンドを実行する", args: []string{"echo", "hello"}, paletteSize: 128, exe: true, w: 5, h: 2, box: []string{"hello", ""}},
			{desc: "マルチバイト文字を含む", args: []string{"あいうえお"}, paletteSize: 128, w: 10, h: 1, box: []string{"BBDDFFHHJJ"}},
			{desc: "シングルバイト・マルチバイト混在", args: []string{"aiueo\nあいうえお"}, paletteSize: 128, w: 10, h: 2, box: []string{"aiueo", "BBDDFFHHJJ"}},
		}
		for _, v := range tds {
			w, h, box := buildString(v.args, v.exe, v.paletteSize)
			if v.w == w && v.h == h && as.Equal(v.box, box) {
				t.Log(fmt.Sprintf("[OK] %s:", v.desc))
			} else {
				var errMsg string
				errMsg += fmt.Sprintf("  w: expect = %d, got = %d\n", v.w, w)
				errMsg += fmt.Sprintf("  h: expect = %d, got = %d\n", v.h, h)
				errMsg += fmt.Sprintf("box: expect = %v, got = %v\n", v.box, box)
				t.Error(fmt.Sprintf("[NG] %s:\n%s", v.desc, errMsg))
			}
		}
	}

	// 引数なし（標準入力）のときのテストケース
	{
		type TestData struct {
			desc        string // テストの目的、説明
			stdin       string
			exe         bool
			paletteSize int
			w, h        int
			box         []string
		}
		tds := []TestData{
			{desc: "標準入力から改行を含まない文字列を受け取る", stdin: "hello", paletteSize: 128, w: 5, h: 1, box: []string{"hello"}},
			{desc: "標準入力から改行を含む文字列を受け取る", stdin: "hello\nworld2", paletteSize: 128, w: 6, h: 2, box: []string{"hello", "world2"}},
			{desc: "シェルのコマンドを実行する", stdin: "echo hello", exe: true, paletteSize: 128, w: 5, h: 2, box: []string{"hello", ""}},
			{desc: "マルチバイト文字を含む", stdin: "あいうえお", paletteSize: 128, w: 10, h: 1, box: []string{"BBDDFFHHJJ"}},
			{desc: "シングルバイト・マルチバイト混在", stdin: "aiueo\nあいうえお", paletteSize: 128, w: 10, h: 2, box: []string{"aiueo", "BBDDFFHHJJ"}},
		}
		for _, v := range tds {
			// 標準入力にテキストを渡す
			stdinRead, stdinWrite, _ := os.Pipe()
			stdinWrite.Write([]byte(v.stdin))
			stdinWrite.Close()
			os.Stdin = stdinRead

			w, h, box := buildString([]string{}, v.exe, v.paletteSize)
			if v.w == w && v.h == h && as.Equal(v.box, box) {
				t.Log(fmt.Sprintf("[OK] %s:", v.desc))
			} else {
				var errMsg string
				errMsg += fmt.Sprintf("  w: expect = %d, got = %d\n", v.w, w)
				errMsg += fmt.Sprintf("  h: expect = %d, got = %d\n", v.h, h)
				errMsg += fmt.Sprintf("box: expect = %v, got = %v\n", v.box, box)
				t.Error(fmt.Sprintf("[NG] %s:\n%s", v.desc, errMsg))
			}
		}
	}
}

func TestFixValueThatOverBoundary(t *testing.T) {
	as := assert.New(t)

	type TestData struct {
		desc     string // テストの目的、説明
		s        string
		boundary int
		expect   string
	}
	tds := []TestData{
		{desc: "アルファベットのみ", s: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", boundary: 128, expect: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{desc: "全角のみ", s: "あいうえお", boundary: 128, expect: "BBDDFFHHJJ"},
	}
	for _, v := range tds {
		got := fixValueThatOverBoundary(v.s, v.boundary)
		as.Equal(v.expect, got, v.desc)
	}
}
