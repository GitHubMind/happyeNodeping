package cmd

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/tealeg/xlsx"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
	"time"
)

const user = "jiangyu.huang@xxx.io"
const pwd = "xxxxxxx"

func GetALlAlert() *cli.Command {
	//Get ALl node ping
	return &cli.Command{
		Name:  "all",
		Usage: "This is command 1",
		Action: func(c *cli.Context) error {
			log.Println("test")
			var err error
			result := make([]FormValues, 0)

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
			ctx, cancel = context.WithTimeout(ctx, 3000*time.Second)
			defer cancel()
			rowCount := 0
			// 替换以下的YOUR_EMAIL和YOUR_PASSWORD为你的登录凭证
			err = chromedp.Run(ctx,
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
			for i := 1; i <= rowCount; i++ {
				formValues := FormValues{}

				xpath := fmt.Sprintf(`	/html/body/div[2]/div[1]/div[2]/div[1]/div[3]/div[1]/div/table/tbody/tr[%d]/td[7]/a/span`, i) // 适当修改以匹配按钮的实际位置
				err := chromedp.Run(ctx,
					chromedp.Text(xpath, &formValues.Result, chromedp.BySearch),
				)
				xpath = fmt.Sprintf(`//*[@id="resultslist"]/tbody/tr[%d]/td[3]/div[2]`, i) // 适当修改以匹配按钮的实际位置
				err = chromedp.Run(ctx,
					chromedp.Click(xpath, chromedp.BySearch),
				)
				if err != nil {
					log.Fatalf("Failed to click button on row %d: %v", i, err)
				}

				time.Sleep(1 * time.Second)
				formValues.notifyMark(ctx)
				value := ""
				err = chromedp.Run(ctx,
					chromedp.WaitVisible(`//*[@id="checksform_label"]`, chromedp.BySearch),
					chromedp.Evaluate(`document.querySelector('#checksform_type').value`, &value),
				)

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
					result = append(result, formValues)
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
					result = append(result, formValues)
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
					result = append(result, formValues)
				}
				if err != nil {
					log.Fatalf("Failed to click button on row %d: %v", i, err)
				}

				err = chromedp.Run(ctx, chromedp.Click(`/html/body/div[3]/div[1]/a/span`, chromedp.BySearch))
				if err != nil {
					log.Println("close", err)
				}
			}
			execl(result, "allnode")
			return nil
		},
	}
}

func execl(value []FormValues, name string) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(name)
	if err != nil {
		fmt.Printf("Error adding sheet: %s\n", err)
		return
	}

	// 添加标题行
	row := sheet.AddRow()
	row.AddCell().Value = "Target"
	row.AddCell().Value = "Location"
	row.AddCell().Value = "Interval(min)"
	row.AddCell().Value = "Result"
	row.AddCell().Value = "Notify Channel"
	row.AddCell().Value = "Tracing API"
	row.AddCell().Value = "SLA"
	row.AddCell().Value = "Target name"

	for _, formValues := range value {
		// 添加一个数据行
		row = sheet.AddRow()
		row.AddCell().Value = formValues.Label
		row.AddCell().Value = formValues.Region
		row.AddCell().Value = formValues.Interval
		row.AddCell().Value = formValues.Result
		row.AddCell().Value = formValues.Notice
		if strings.Contains("state_traceBlock", formValues.PostData) {
			row.AddCell().Value = "Yes"
		} else {
			row.AddCell().Value = "No"
		}
		row.AddCell().Value = ""
		row.AddCell().Value = formValues.Type
		row.AddCell().Value = formValues.TargetURL
		row.AddCell().Value = formValues.PostData

	}

	// 保存文件
	err = file.Save(name + ".xlsx")
	if err != nil {
		fmt.Printf("Error saving file: %s\n", err)
	}
}

func (fv *FormValues) notifyMark(ctx context.Context) {
	notifyNumber := 0
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelectorAll("#notificationslist tbody tr").length`, &notifyNumber),
	)
	if err != nil {
		log.Println("rows notice", err)
	}
	notifyArr := make([]string, notifyNumber-1, notifyNumber-1)

	for j := 0; j < notifyNumber-1; j++ {
		notify := ""
		///html/body/div[3]/div[2]/div/div[3]/div/table/tbody/tr[2]/td[1]/text()
		path := fmt.Sprintf(`document.evaluate('/html/body/div[3]/div[2]/div/div[3]/div/table/tbody/tr[%d]/td[1]/text()', document, null, XPathResult.ANY_TYPE, null).iterateNext().textContent`, (j + 2))
		//log.Println(path)
		err = chromedp.Run(ctx,
			chromedp.Evaluate(path, &notify),
		)
		notifyArr = append(notifyArr, notify)
		if err != nil {

			log.Println("err", err)
		}
	}
	if len(notifyArr) > 0 {
		fv.Notice = strings.Join(notifyArr, "/n")
	}
}
