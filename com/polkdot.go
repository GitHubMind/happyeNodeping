package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type FormValues struct {
	Type      string `json:"type"`
	Interval  string `json:"interval"`
	TargetURL string `json:"target_url"`
	PostData  string `json:"param_postdata"`
	Label     string `json:"checksform_label"`
	Region    string `json:"checksform_runlocations"`
	Result    string `json:"result"`
	Notice    string `json:"notice"`
}

func Node() *cli.Command {
	// 读取 polkdot 配置 ,然后根据网络去判断那些是高可用节点.
	// 这个是处理 collection list
	return &cli.Command{
		Name:  "node",
		Usage: "This is command 1",
		Action: func(c *cli.Context) error {

			log.Println("Running command 1")
			//这个文件是nancy 整理的,需要感谢下nancy肉眼去收集
			jsonFile, err := os.Open("network.json")
			if err != nil {
				panic(err)
			}
			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)

			var networkMap map[string]string
			json.Unmarshal(byteValue, &networkMap)

			srcFile, err := os.Open("sourceNode.xlsx")
			if err != nil {
				panic(err)
			}
			defer srcFile.Close()

			destFile, err := os.Create("resultNode.csv")
			if err != nil {
				panic(err)
			}
			defer destFile.Close()

			reader := csv.NewReader(srcFile)

			// 读取源CSV文件的每一行，并添加一个新字段
			resultReader := make([][]string, 0, 0)
			polkdotIndex := -1
			networkIndex := -1
			ApiIndex := -1
			nonSlaIndex := -1

			apiNumber := make(map[string]int)
			apiName := make(map[string]map[string]int)
			for {
				record, err := reader.Read()
				if err != nil {
					break
				}

				resultReader = append(resultReader, record)
				if polkdotIndex == -1 {
					for i, s := range record {
						if s == "polkadot.js" {
							polkdotIndex = i
						}
					}
				}
				//通过 excel的 第一排名字去作为依据判断
				if networkIndex == -1 {
					for i, s := range record {
						if strings.TrimSpace(s) == "Network" {
							networkIndex = i
						}
					}
				}
				if nonSlaIndex == -1 {
					for i, s := range record {

						if strings.TrimSpace(s) == "non-sla" {
							nonSlaIndex = i
						}
					}
				}
				if ApiIndex == -1 {
					for i, s := range record {
						if strings.TrimSpace(s) == "Service ID" {
							ApiIndex = i
						}
					}
				} else {
					if apiName[record[ApiIndex]] == nil {
						apiName[record[ApiIndex]] = make(map[string]int)
					}
					//polkadex-xx-x-txxxng			Polkadex	这是一堆数字 	polkadex	running	t	archive	110Gi	POLKADOT & PARACHAINS
					apiName[record[ApiIndex]][record[0]]++
					if apiName[record[ApiIndex]][record[0]] == 1 {
						apiNumber[record[ApiIndex]]++
					}

				}
				// 写入新的CSV文件

			}

			writer := csv.NewWriter(destFile)

			for i, i2 := range resultReader {

				val, ok := networkMap[i2[networkIndex]]
				if i > 0 && ok {
					i2[polkdotIndex] = val
				}
				valInt, ok := apiNumber[i2[ApiIndex]]
				if i > 0 && valInt == 1 {
					i2[nonSlaIndex] = "true"

				}

				if err := writer.Write(i2); err != nil {
					panic(err)
				}
			}

			writer.Flush()
			if err := writer.Error(); err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}
}
