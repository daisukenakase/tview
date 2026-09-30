package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/muesli/termenv"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"golang.org/x/term"
)

const asciiChars = " .:-=+*#%@"

func main() {
	colorMode := flag.Bool("color", false, "カラー表示モード(色ブロック)")
	colorAsciiMode := flag.Bool("color-ascii", false, "カラー表示モード(色付き文字)")
	targetWidth := flag.Int("width", 0, "出力する最大横幅(省略時は端末幅に合わせる)")
	flag.Parse()
	modeSpecified := false
	flag.Visit(func(argument *flag.Flag) {
		if argument.Name == "color" || argument.Name == "color-ascii" {
			modeSpecified = true
		}
	})
	if !modeSpecified {
		*colorMode = true
	}
	*colorMode, *colorAsciiMode = applyColorFallback(*colorMode, *colorAsciiMode, termenv.NewOutput(os.Stdout).ColorProfile())
	output := colorable.NewColorable(os.Stdout)

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("使用方法: ./tview [-color | -color-ascii] [-width 80] <画像パス>")
		os.Exit(1)
	}
	imagePath := args[0]

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: ファイルを開けません: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 画像をデコードできません: %v\n", err)
		os.Exit(1)
	}

	bounds := img.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	maxWidth := *targetWidth
	if maxWidth < 1 {
		maxWidth = 80
	}
	maxHeight := 0
	if columns, rows, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		if terminalWidth := columns / 2; terminalWidth > 0 && (*targetWidth < 1 || terminalWidth < maxWidth) {
			maxWidth = terminalWidth
		}
		if rows > 0 {
			maxHeight = rows
		}
	}
	w, h := fitDimensions(origW, origH, maxWidth, maxHeight)

	resizedImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.BiLinear.Scale(resizedImg, resizedImg.Bounds(), img, bounds, draw.Over, nil)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := resizedImg.At(x, y).RGBA()
			R := int(r >> 8)
			G := int(g >> 8)
			B := int(b >> 8)

			switch {
			case *colorMode:
				fmt.Fprintf(output, "\x1b[48;2;%d;%d;%dm  ", R, G, B)
			case *colorAsciiMode:
				gray := 0.299*float64(R) + 0.587*float64(G) + 0.114*float64(B)
				charIdx := int(gray / 256.0 * float64(len(asciiChars)))
				if charIdx >= len(asciiChars) {
					charIdx = len(asciiChars) - 1
				}
				char := string(asciiChars[charIdx])
				fmt.Fprintf(output, "\x1b[38;2;%d;%d;%dm%s", R, G, B, char+char)
			default:
				gray := 0.299*float64(R) + 0.587*float64(G) + 0.114*float64(B)
				charIdx := int(gray / 256.0 * float64(len(asciiChars)))
				if charIdx >= len(asciiChars) {
					charIdx = len(asciiChars) - 1
				}
				char := string(asciiChars[charIdx])
				fmt.Fprint(output, char+char)
			}
		}
		if *colorMode || *colorAsciiMode {
			fmt.Fprint(output, "\x1b[0m")
		}
		fmt.Fprintln(output)
	}
}

func applyColorFallback(colorMode, colorAsciiMode bool, profile termenv.Profile) (bool, bool) {
	if profile == termenv.Ascii {
		return false, false
	}
	return colorMode, colorAsciiMode
}

func fitDimensions(origW, origH, maxWidth, maxHeight int) (int, int) {
	width := maxWidth
	if maxHeight > 0 {
		widthForHeight := int(float64(maxHeight)*float64(origW)/float64(origH) + 0.5)
		if widthForHeight < width {
			width = widthForHeight
		}
	}
	if width < 1 {
		width = 1
	}

	height := int(float64(origH)/float64(origW)*float64(width) + 0.5)
	if height < 1 {
		height = 1
	}
	if maxHeight > 0 && height > maxHeight {
		height = maxHeight
	}
	return width, height
}
