package models

/**
 * @File: response.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/25 上午10:33
 * @Software: GoLand
 * @Version:  1.0
 */

type JsonResult struct {
	StatusCode string      `json:"status_code"`
	Msg        interface{} `json:"msg"`
	MsgDetail  string      `json:"msg_detail"`
	Data       interface{} `json:"data"`
}
