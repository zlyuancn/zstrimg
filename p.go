/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/12/31
   Description :
-------------------------------------------------
*/

package zstrimg

import (
    "bytes"
    "image"
    "image/color"
    "image/draw"
    "image/jpeg"
    "image/png"
    "io"
    "os"
    "strings"

    "github.com/zlyuancn/zstrimg/interp"
)

var _Chr = []byte(" .,-:;irs1h5S398GX&AHBM#@")

func scale(src image.Image, out draw.Image) {
    b := out.Bounds()
    srcb := src.Bounds()
    sx := float64(b.Dx()) / float64(srcb.Dx())
    sy := float64(b.Dy()) / float64(srcb.Dy())
    _ = I.Scale(sx, sy).Transform(out, src, interp.Bilinear)
}

// 按宽度等比缩放
func ScaleImage(src image.Image, newdx int) *image.RGBA {
    bound := src.Bounds()
    dx := bound.Dx()
    dy := bound.Dy()
    dst := image.NewRGBA(image.Rect(0, 0, newdx, newdx*dy/dx))
    scale(src, dst)
    return dst
}

// 指定宽高缩放
func ScaleImageEx(src image.Image, newdx, newdy int) *image.RGBA {
    dst := image.NewRGBA(image.Rect(0, 0, newdx, newdy))
    scale(src, dst)
    return dst
}

// 将编码对象进行处理后返回字节数组
func SaveImage(img *image.RGBA, filetype string, w io.Writer) (err error) {
    switch filetype {
    case "png":
        err = png.Encode(w, img)
    case "jpg", "jpeg":
        err = jpeg.Encode(w, img, nil)
    default:
        err = jpeg.Encode(w, img, nil)
    }
    return
}

// 图片转灰度
func HDImage(m image.Image) *image.RGBA {
    bounds := m.Bounds()
    dx := bounds.Dx()
    dy := bounds.Dy()
    newRgba := image.NewRGBA(bounds)
    for i := 0; i < dx; i++ {
        for j := 0; j < dy; j++ {
            colorRgb := m.At(i, j)
            _, g, _, a := colorRgb.RGBA()
            g_uint8 := uint8(g >> 8)
            a_uint8 := uint8(a >> 8)
            newRgba.SetRGBA(i, j, color.RGBA{g_uint8, g_uint8, g_uint8, a_uint8})
        }
    }
    return newRgba
}

func colorToString(c int) byte {
    c = int(float32(c) / 255 * float32(len(_Chr)-1))
    return _Chr[c]
}

// 图片转字符
func ImageToString(m image.Image, sep string) string {
    bounds := m.Bounds()
    dx := bounds.Dx()
    dy := bounds.Dy()

    out := make([]string, dy)

    for j := 0; j < dy; j++ {
        line := make([]byte, dx)
        for i := 0; i < dx; i++ {
            colorRgb := m.At(i, j)
            _, g, _, _ := colorRgb.RGBA()
            line[i] = colorToString(int(g >> 8))
        }
        out[j] = string(line)
    }
    return strings.Join(out, sep)
}

// 从文件加载图片
func LoadImageOfFile(filename string) (image.Image, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    return LoadImageOfReader(f)
}

// 从读取器中加载图片
func LoadImageOfReader(r io.Reader) (image.Image, error) {
    img, _, err := image.Decode(r)
    return img, err
}

// 从byte中加载图片
func LoadImageOfByte(data []byte) (image.Image, error) {
    return LoadImageOfReader(bytes.NewReader(data))
}

// 文件图片转字符图
func ImageFileToString(filename string, newdx, newdy int, sep string) (string, error) {
    img, err := LoadImageOfFile(filename)
    if err != nil {
        return "", err
    }

    out := ScaleImageEx(img, newdx, newdy)
    out = HDImage(out)
    return ImageToString(out, sep), nil
}
