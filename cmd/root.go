// Copyright © 2019 xztaityozx
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xztaityozx/t2p/nutils"
	"github.com/xztaityozx/t2p/palette"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "t2p",
	Short: "",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		exe, _ := cmd.Flags().GetBool("execute")
		height, _ := cmd.Flags().GetInt("height")
		width, _ := cmd.Flags().GetInt("width")
		size, _ := cmd.Flags().GetInt("size")
		table, _ := cmd.Flags().GetString("table")

		//
		w, h, box := buildString(args, exe)

		// フラグの指定値が優先
		if height == 0 {
			height = h
		}

		if width == 0 {
			width = w
		}

		p := palette.NewPalette(table)
		img := p.Create(width, height, box)

		if err := outImage(out, img, size); err != nil {
			logrus.WithError(err).Fatal("Failed encode image")
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("out", "o", "", "画像の出力先です。defaultはstdoutです")
	rootCmd.Flags().BoolP("execute", "e", false, "入力文字列をコマンドとして受け取り実行結果を画像にします")
	rootCmd.Flags().IntP("width", "W", 0, "出力画像の幅です。defaultは文字列の長さです")
	rootCmd.Flags().IntP("height", "H", 0, "出力画像の高さです。defaultは1ドットです")
	rootCmd.Flags().IntP("size", "s", 100, "出力画像の拡大率(%)です")
	rootCmd.Flags().String("table", "", "変換テーブルを指定できます。合計127ドットのpngファイルが指定できます")
}

func buildString(args []string, exe bool) (int, int, []string) {
	str := ""
	if len(args) == 0 {
		s := bufio.NewScanner(os.Stdin)
		var b []string
		for s.Scan() {
			b = append(b, s.Text())
		}
		str = strings.Join(b, "\n")
	} else {
		str = strings.Join(args, " ")
	}

	if exe {
		command := str
		out, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			logrus.WithError(err).Fatal("Failed execute command: ", command)
		}
		str = string(out)
	}

	l := strings.Split(str, "\n")
	m := 0
	for _, v := range l {
		m = nutils.IntMax(m, len(v))
	}
	return m, len(l), l
}

func outImage(path string, src *image.RGBA, size int) error {

	var fp *os.File
	var format = ".png"
	if len(path) == 0 {
		fp = os.Stdout
	} else {
		var err error
		fp, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			logrus.WithError(err).Fatal("Failed open or create file: ", path)
		}

		format = filepath.Ext(path)
	}
	defer fp.Close()

	img := palette.NewZoom(size).ScaleUp(src)

	if format == ".png" {
		return png.Encode(fp, img)
	} else if format == ".jpg" {
		return jpeg.Encode(fp, img, nil)
	} else if format == ".gif" {
		return gif.Encode(fp, img, nil)
	} else {
		logrus.Fatal("t2p does not support file type: ", format)
	}

	return nil
}
