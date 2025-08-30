package component

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func SenEamil(duration string, to ...string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "2908451607@qq.com")
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "服务运行状态")
	m.SetBody("text/html", storeHTML(duration))
	d := gomail.NewDialer("smtp.qq.com", 465, "2908451607@qq.com", "lfloolaeismedcfi")
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("发送成功")
}
func storeHTML(duration string) string {
	return "<!DOCTYPE html>\n<html>\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>服务运行状态通知</title>\n    <style type=\"text/css\">\n        body {\n            font-family: Arial, sans-serif;\n            line-height: 1.6;\n            color: #333333;\n            background-color: #f5f5f5;\n            margin: 0;\n            padding: 0;\n        }\n        .email-container {\n            max-width: 600px;\n            margin: 20px auto;\n            background-color: #ffffff;\n            border-radius: 8px;\n            overflow: hidden;\n            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);\n        }\n        .header {\n            background-color: #4CAF50;\n            color: white;\n            padding: 20px;\n            text-align: center;\n        }\n        .header h1 {\n            margin: 0;\n            font-size: 24px;\n        }\n        .content {\n            padding: 30px;\n        }\n        .status {\n            text-align: center;\n            margin: 20px 0;\n            padding: 15px;\n            border-radius: 4px;\n            font-weight: bold;\n        }\n        .running {\n            background-color: #E8F5E9;\n            color: #2E7D32;\n            border-left: 4px solid #4CAF50;\n        }\n        .footer {\n            text-align: center;\n            padding: 15px;\n            background-color: #f9f9f9;\n            font-size: 12px;\n            color: #777777;\n        }\n        .button {\n            display: inline-block;\n            padding: 10px 20px;\n            background-color: #4CAF50;\n            color: white;\n            text-decoration: none;\n            border-radius: 4px;\n            margin: 20px 0;\n        }\n        .info-table {\n            width: 100%;\n            border-collapse: collapse;\n            margin: 20px 0;\n        }\n        .info-table th, .info-table td {\n            padding: 12px;\n            text-align: left;\n            border-bottom: 1px solid #dddddd;\n        }\n        .info-table th {\n            background-color: #f2f2f2;\n        }\n    </style>\n</head>\n<body>\n    <div class=\"email-container\">\n        <div class=\"header\">\n            <h1>服务运行状态通知</h1>\n        </div>\n        \n        <div class=\"content\">\n            <p>您好，</p>\n            \n            <p>这是来自您服务的定时状态通知，当前服务运行正常。</p>\n            \n            <div class=\"status running\">\n                <p>状态: 运行中 ✅</p>\n                <p>最后检查时间: <span id=\"current-time\">" + NowDateTime() + "</span></p>\n            </div>\n            \n            <h3>服务信息</h3>\n            <table class=\"info-table\">\n                <tr>\n                    <th>服务名称</th>\n                    <td>golang爬虫</td>\n                </tr>\n                <tr>\n                    <th>服务器</th>\n                    <td>8.138.180.182</td>\n                </tr>\n                <tr>\n                    <th>运行时间</th>\n                    " +
		fmt.Sprintf("<td>%s</td>\n", duration) +
		"                </tr>\n                <!-- <tr>\n                    <th>CPU使用率</th>\n                    <td>23%</td>\n                </tr>\n                <tr>\n                    <th>内存使用率</th>\n                    <td>45%</td>\n                </tr> -->\n            </table>\n            \n            <!-- <p>如果您需要立即检查服务状态，请点击下方按钮：</p> -->\n            <!-- <a href=\"https://your-service-dashboard.com\" class=\"button\">查看服务仪表盘</a> -->\n            \n            <p>此邮件为系统自动发送，请勿直接回复。</p>\n        </div>\n        \n        <div class=\"footer\">\n            <p>© 2023 gxuicpc 保留所有权利.</p>\n        </div>\n    </div>\n    \n    <!-- 使用JavaScript动态设置当前时间（实际邮件中不会执行，仅用于预览） -->\n    <script>\n        </script>\n</body>\n</html>"
}
