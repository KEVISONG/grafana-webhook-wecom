package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/KEVISONG/go/common/http"
	"github.com/grafana/grafana/pkg/services/alerting"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/sirupsen/logrus"
)

var (
	webhook string
	port    int
	err     error
)

type WebhookBody struct {
	Title       string
	RuleID      string
	Message     string
	ImageURL    string
	EvalMatches []*alerting.EvalMatch
}

// PostNotification PostNotification
func PostNotification(ctx iris.Context) {

	/*
		evalContext := &alerting.EvalContext{}

		ctx.ReadJSON(evalContext)

		bytes, err := json.Marshal(*evalContext)
		if err != nil {
			logrus.Error(err)
		}
		fmt.Println(string(bytes))
	*/
	webhookBody := &WebhookBody{}
	ctx.ReadJSON(webhookBody)
	content := fmt.Sprintf("%v\n%v\n",
		webhookBody.Title,
		webhookBody.Message,
	)

	if webhookBody.ImageURL != "" {
		content += fmt.Sprintf("[%s](%s)\n", webhookBody.ImageURL, webhookBody.ImageURL)
	}

	for i, match := range webhookBody.EvalMatches {
		content += fmt.Sprintf("\n%2d. %s: `%s`", i+1, match.Metric, match.Value)
	}

	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": content,
		},
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		logrus.Error("Failed to marshal body", "error", err)
		return
	}

	resp, err := http.SetHeaders(map[string]string{"content-type": "application/json"}).Post(webhook, bodyJSON)
	if err != nil {
		logrus.Error(err)
		return
	}

	fmt.Println(string(resp))

}

func serve(port int) {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Post("/api/webhook/wecom", PostNotification)
	app.Run(iris.Addr(fmt.Sprintf("0.0.0.0:%v", port)))
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify WeCom Robot API")
		os.Exit(1)
	}

	if len(os.Args) >= 2 {
		webhook = os.Args[1]
	}

	if len(os.Args) >= 3 {
		port, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Port must be positive")
			os.Exit(1)
		}
	} else {
		port = 80
	}

	serve(port)

}
