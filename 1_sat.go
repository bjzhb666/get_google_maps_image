package main

import (
	"encoding/json"

	"./lib"
	"github.com/ironsublimate/gomapinfer/common"
	"github.com/ironsublimate/gomapinfer/googlemaps"

	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

const ZOOM = 20

func main() {
	apiKey := os.Args[1]
	outDir := os.Args[2]

	widthWorld := 2 * math.Pi * 6378137 / math.Exp2(ZOOM) / 256 // meters per pixel
	regions := lib.GetRegions()

	type Tile struct {
		Region   lib.Region
		X        int
		Y        int
		Filename string
	}

	type Result struct {
		CenterGPS   string `json:"centerGPS"`
		CenterWorld string `json:"centerWorld"`
		Filename    string `json:"filename"`
	}
	// 创建一个结果切片
	results := make([]Result, 0)

	fmt.Printf("found %d regions\n", len(regions))
	fmt.Printf("widthWorld: %f\n", widthWorld)
	fmt.Printf(("apiKey: %s\n"), apiKey)
	// fmt.Printf(regions)

	var requiredTiles []Tile
	for _, region := range regions {
		for x := -region.RadiusX; x < region.RadiusX; x++ {
			for y := -region.RadiusY; y < region.RadiusY; y++ {
				fname := fmt.Sprintf("%s/%s_%d_%d_sat.png", outDir, region.Name, x, y)
				if _, err := os.Stat(fname); err == nil {
					continue
				}
				requiredTiles = append(requiredTiles, Tile{
					Region:   region,
					X:        x,
					Y:        y,
					Filename: fname,
				})
			}
		}
	}

	fmt.Printf("found %d required tiles\n", len(requiredTiles)) // if the file is already there, it will not be created
	// fmt.Println(requiredTiles)

	for _, tile := range requiredTiles {
		fmt.Printf("creating %s\n", filepath.Base(tile.Filename))
		im := image.NewNRGBA(image.Rect(0, 0, 4096, 4096))
		for xOffset := 0; xOffset < 4096; xOffset += 512 {
			for yOffset := 0; yOffset < 4096; yOffset += 512 {
				centerWorld := tile.Region.CenterWorld.Add(common.Point{float64(tile.X*4096 + xOffset),
					float64(-(tile.Y*4096 + yOffset))}.Scale(widthWorld))
				// fmt.Println("centerWorld", centerWorld)
				// fmt.Println((tile.X*4096 + xOffset), -(tile.Y*4096 + yOffset)) // -512, 512
				centerGPS := googlemaps.MetersToLonLat(centerWorld)
				// fmt.Println("centerGPS", centerGPS)
				result := Result{
					CenterGPS:   fmt.Sprintf("%v", centerGPS),
					CenterWorld: fmt.Sprintf("%v", centerWorld),
					Filename:    tile.Filename,
				}
				results = append(results, result)

				satelliteImage := googlemaps.GetSatelliteImage(centerGPS, ZOOM, apiKey)
				for i := 0; i < 512; i++ {
					for j := 0; j < 512; j++ {
						im.Set(xOffset+i, yOffset+j, satelliteImage.At(i, j))
					}
				}
			}
		}
		f, err := os.Create(tile.Filename)
		if err != nil {
			panic(err)
		}
		if err := png.Encode(f, im); err != nil {
			panic(err)
		}
		f.Close()
	}
	// fmt.Println("results", results)
	//将results切片转换为map
	resultsMap := make(map[string][]Result)
	for _, result := range results {
		// 将结果添加到对应的切片中
		resultsMap[result.Filename] = append(resultsMap[result.Filename], result)
	}
	//将map转换为JSON
	jsonData, err := json.Marshal(resultsMap)
	if err != nil {
		panic(err)
	}
	//创建JSON文件
	jsonFile, err := os.Create(outDir + "/results.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	//将JSON写入文件
	jsonFile.WriteString(string(jsonData))

	// resultsMap := make(map[string]Result)
	// for _, result := range results {
	// 	resultsMap[result.Filename] = result
	// }

	// // 将map转换为JSON
	// jsonData, err := json.Marshal(resultsMap)
	// if err != nil {
	// 	panic(err)
	// }

	// // 创建JSON文件
	// jsonFile, err := os.Create(outDir + "/results.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer jsonFile.Close()

	// // 将JSON写入文件
	// jsonFile.WriteString(string(jsonData))
}
