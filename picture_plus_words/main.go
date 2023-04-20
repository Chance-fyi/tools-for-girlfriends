package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"os"
	"time"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	templatePicture := viper.GetString("templatePicture")
	img, err := gg.LoadPNG(templatePicture)
	if err != nil {
		fmt.Println(templatePicture+" 图片不存在", err.Error())
		return
	}
	w := img.Bounds().Size().X
	h := img.Bounds().Size().Y

	dc := gg.NewContext(w, h)

	err = dc.LoadFontFace(viper.GetString("font"), viper.GetFloat64("fontSize"))
	if err != nil {
		fmt.Println("字体加载失败", err.Error())
		return
	}
	dc.SetRGB(viper.GetFloat64("fontColorR"), viper.GetFloat64("fontColorG"), viper.GetFloat64("fontColorB"))

	dir := "./" + time.Now().Format("2006-01-02")
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 777)
		if err != nil {
			fmt.Println("目录创建失败", err.Error())
			return
		}
	}

	for _, row := range getName() {
		if len(row) < 2 {
			continue
		}
		id, name := row[0], row[1]
		if name == "" {
			continue
		}
		dc.Clear()
		dc.DrawImage(img, 0, 0)
		dc.DrawStringAnchored(viper.GetString("prefix")+name+viper.GetString("suffix"), viper.GetFloat64("x"), viper.GetFloat64("y"), viper.GetFloat64("ax"), viper.GetFloat64("ay"))
		err = dc.SavePNG(fmt.Sprintf("./%s/%s_%s.png", dir, id, name))
		if err != nil {
			fmt.Println(name+" 生成失败", err.Error())
		}
	}
}

func getName() [][]string {
	f, err := excelize.OpenFile("name.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rows[1:]
}
