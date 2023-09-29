package cmd

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/tealeg/xlsx"
	"github.com/urfave/cli/v2"
	"log"
	"time"
)

var locationMap = map[string]string{}

// TODO 优化输入路径
var filePath = "./Alerting.xlsx"

func UpdsteALlAlert() *cli.Command {
	//Get ALl node ping
	return &cli.Command{
		Name:  "update",
		Usage: "update label",
		Action: func(c *cli.Context) error {

			//这是每个读取的对应
			locationMap = make(map[string]string, 0)
			locationMap["EU"] = "eur"
			locationMap["AP"] = "eao"
			locationMap["LN"] = "lam"
			locationMap["World"] = "wlw"
			locationMap["NA"] = "nam"
			// 替换名字的excel表
			strMap := check()
			// 设置ExecAllocatorOptions
			options := append(chromedp.DefaultExecAllocatorOptions[:],
				// 禁用headless模式
				chromedp.Flag("headless", false),
				chromedp.Flag("disable-gpu", true), // 当需要显示GUI时，通常也推荐这个选项
			)

			// 创建新的chromedp上下文
			ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
			defer cancel()

			ctx, cancel = chromedp.NewContext(ctx)
			defer cancel()

			// 设置超时
			ctx, cancel = context.WithTimeout(ctx, 1000*time.Second)
			defer cancel()
			rowCount := 0
			// 替换以下的YOUR_EMAIL和YOUR_PASSWORD为你的登录凭证
			err := chromedp.Run(ctx,
				chromedp.Navigate(`https://nodeping.com/`),
				chromedp.Click(`/html/body/header/div/div[4]/ul/li[1]/a`, chromedp.BySearch),

				chromedp.SendKeys(`/html/body/section/div/div[2]/form/ul/li[1]/input`, user, chromedp.BySearch),
				chromedp.SendKeys(`/html/body/section/div/div[2]/form/ul/li[2]/input`, pwd, chromedp.BySearch),
				chromedp.Click(`/html/body/section/div/div[2]/form/ul/li[3]/input`, chromedp.BySearch),
				chromedp.WaitVisible(`/html/body/div[2]/div[1]/div[2]/div[1]/div[3]/div[1]/div/div[1]`, chromedp.BySearch),
				chromedp.WaitVisible(`#resultslist_length`, chromedp.ByID), // 等待下拉列表元素出现
				chromedp.SetValue(`//*[@id="resultslist_length"]/label/select`, "-1", chromedp.BySearch),
				chromedp.Evaluate(`document.querySelectorAll("#resultslist tbody tr").length`, &rowCount),
			)
			fmt.Printf("Row count: %d\n", rowCount)
			if err != nil {
				log.Fatalf("Failed running task: %v", err)
			}
			//处理弹窗
			go func() {
				for true {
					time.Sleep(1 * time.Second)
					ctx2, cancel1 := context.WithTimeout(context.TODO(), 1*time.Second)
					defer cancel1()
					var nodeInfo []*cdp.Node
					xpath := `//*[@id="np_messagedialog"]`
					err := chromedp.Run(ctx2,
						chromedp.Nodes(xpath, &nodeInfo, chromedp.BySearch),
					)
					if err != nil {
						log.Println("Error during chromedp Run:", err)
						continue
						// 处理错误D
					} else if nodeInfo != nil {
						log.Println("Element exists!")
					} else {
						log.Println("Element  no exists!")
						continue
					}

					err = chromedp.Run(ctx,
						chromedp.Click(`/html/body/div[3]/div[11]/div/button/span`, chromedp.NodeVisible),
					)
					if err != nil {
						log.Println("click", err)
						continue
					}

				}
			}()

			for i := 1; i <= rowCount; i++ {
				log.Println(i)
				ctx1, cancel1 := context.WithTimeout(ctx, 10*time.Second)
				defer cancel1()
				formValues := FormValues{}

				xpath := fmt.Sprintf(`/html/body/div[2]/div[1]/div[2]/div[1]/div[3]/div[1]/div/table/tbody/tr[%d]/td[3]/div[1]/a`, i) // 适当修改以匹配按钮的实际位置
				err = chromedp.Run(ctx1,
					chromedp.Text(xpath, &formValues.Result, chromedp.BySearch),
				)
				if err != nil {
					time.Sleep(1 * time.Second)
					log.Println(err)
					continue
					//log.Fatalf("Failed to click button on row %d: %v", i, err)
				}
				if _, ok := strMap[formValues.Result]; ok {
					time.Sleep(3 * time.Second)
				} else {
					continue
				}
				//*[@id="resultslist"]/tbody/tr[540]/td[3]/div[2]

				xpath = fmt.Sprintf(`/html/body/div[2]/div[1]/div[2]/div[1]/div[3]/div[1]/div/table/tbody/tr[%d]/td[3]/div[2]`, i) // 适当修改以匹配按钮的实际位置
				log.Println(xpath)
				err = chromedp.Run(ctx1,
					chromedp.WaitVisible(xpath, chromedp.BySearch),
					chromedp.Click(xpath, chromedp.BySearch),
				)
				if err != nil {
					log.Println(err)
					i--
					continue

					//log.Fatalf("Failed to click button on row %d: %v", i, err)
				}
				err = chromedp.Run(ctx1,
					chromedp.WaitVisible(`//*[@id="checksform_label"]`, chromedp.BySearch), //
					chromedp.Evaluate(`document.querySelector('#checksform_label').value`, &formValues.Label),
				)
				if err != nil {
					i--
					continue

					//log.Fatalf("Failed to click button on row %d: %v", i, err)
				}
				time.Sleep(2 * time.Second)
				err = chromedp.Run(ctx,
					chromedp.SetValue(`//*[@id="checksform_label"]`, strMap[formValues.Result], chromedp.BySearch),
				)
				if err != nil {
					i--
					continue
					//log.Fatalf("Failed to click button on row %d: %v", i, err)
				}

				value := ""
				err = chromedp.Run(ctx,
					chromedp.WaitVisible(`//*[@id="checksform_label"]`, chromedp.BySearch),
					chromedp.Evaluate(`document.querySelector('#checksform_type').value`, &value),
				)
				// TODO 需要补充的类型
				switch value {
				case "SSL":
					log.Println("SSL!!!!!")
					err = chromedp.Run(ctx,
						chromedp.Evaluate(`document.querySelector('#checksform_type').value`, &formValues.Type),
						chromedp.Evaluate(`document.querySelector('#checksform_interval').value`, &formValues.Interval),
						chromedp.Evaluate(`document.querySelector("#checkparams > div:nth-child(1) > input").value`, &formValues.TargetURL),
						chromedp.Evaluate(`document.querySelector("#checksform_label").value`, &formValues.Label),
						chromedp.Evaluate(`document.querySelector("#checksform_runlocations").value`, &formValues.Region),
					)
				case "HTTPADV":
					log.Println("HTTPADV!!!!!")
					err = chromedp.Run(ctx,
						chromedp.Evaluate(`document.querySelector('#checksform_type').value`, &formValues.Type),
						chromedp.Evaluate(`document.querySelector('#checksform_interval').value`, &formValues.Interval),
						chromedp.Evaluate(`document.querySelector('#target_url').value`, &formValues.TargetURL),
						chromedp.Evaluate(`document.querySelector('#param_postdata').value`, &formValues.PostData),
						chromedp.Evaluate(`document.querySelector('#checksform_label').value`, &formValues.Label),
						chromedp.Evaluate(`document.querySelector("#checksform_runlocations").value`, &formValues.Region),
					)
				case "WEBSOCKET":
					log.Println("WEBSOCKET!!!!!")
					err = chromedp.Run(ctx,
						chromedp.Evaluate(`document.querySelector('#checksform_type').value`, &formValues.Type),
						chromedp.Evaluate(`document.querySelector('#checksform_interval').value`, &formValues.Interval),
						chromedp.Evaluate(`document.querySelector('#target_url').value`, &formValues.TargetURL),
						chromedp.Evaluate(`document.querySelector('#data').value`, &formValues.PostData),
						chromedp.Evaluate(`document.querySelector('#checksform_label').value`, &formValues.Label),
						chromedp.Evaluate(`document.querySelector("#checksform_runlocations").value`, &formValues.Region),
					)

				}

				// TODO 这是排除 同一个label(title) 有不同的条件 可以通过以下的方法去排除
				//pass := false
				//for s, s2 := range locationMap {
				//	if strings.Contains(formValues.Label, s) && formValues.Region == s2 {
				//		pass = true
				//		break
				//	}
				//}
				//
				//if !pass {
				//	//获取完信息然后把她给排除掉
				//	log.Println("close", formValues.Label)
				//	err = chromedp.Run(ctx, chromedp.Click(`/html/body/div[3]/div[1]/a/span`, chromedp.BySearch))
				//	continue
				//}

				log.Println(formValues.Result, "->", strMap[formValues.Result])
				//这是获取保存的按钮
				err = chromedp.Run(ctx, chromedp.Click(`/html/body/div[3]/div[11]/div/button[1]/span`, chromedp.BySearch))

				if err != nil {
					i--
					log.Println("close", err)
				}

			}
			return nil
		},
	}

}

// 读取要修改的文件
func check() map[string]string {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	strMap := make(map[string]string)

	for _, sheet := range xlFile.Sheets {
		for key, row := range sheet.Rows {
			if key > 0 && len(row.Cells) > 1 && row.Cells[0].String() != "" && row.Cells[1].String() != "" { // 确保行中至少有两个单元格
				target := row.Cells[0].String()
				changeToName := row.Cells[1].String()
				strMap[target] = changeToName
			}

		}
	}
	return strMap
}
