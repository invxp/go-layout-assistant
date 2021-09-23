package assistant

import "time"

//cronAlert 接收HTTPPost请求后通过CRON进行输出(示例)
func (assistant *Assistant) cronAlert() {
	assistant.logf("cron alert: %s", time.Now().Format("2006-01-02 15:04:05"))
}
