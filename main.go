package main

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"
)

const tip4 = `请选择：
- 1 常开
- 2 常闭
- 3 默认`

var rec []byte

func main() {
	var ip string
	var door2 string
	tip1 := "请输入IP地址（例如：192.168.1.1:60000）:"
	fmt.Println(tip1)
	_, _ = fmt.Scanln(&ip)
	tip2 := "请输入门ID（例如ABCD，只需前四位）："
	fmt.Println(tip2)
	_, _ = fmt.Scanln(&door2)
	door21, _ := strconv.ParseInt(door2[:2], 16, 32)
	door22, _ := strconv.ParseInt(door2[2:], 16, 32)
	door := []byte{uint8(door21), uint8(door22)}
	fmt.Println(door)
	buff := make([]byte, 34)
	udpAddr, _ := net.ResolveUDPAddr("udp4", ip)
	udpConn, _ := net.DialUDP("udp", nil, udpAddr)
	go func() {
		rec = make([]byte, 1024)
		_, _ = udpConn.Read(rec)
		fmt.Println("接收数据:", string(rec))
		if rec[0] == 1 {
			fmt.Println("操作成功！")
		}
	}()
	var control []byte
	var mode string
	const tip3 = `功能列表
- open_door 开门
- time 时间同步
- logcount 获取日志条数
- resetwarn 警报复位
- setdoor 设置门的状态（常开/常闭/默认）
- clear 清空日志
- readcard 获取卡号
- readall 批量获取卡号
- addpriv 添加权限
- rempriv 清除权限
- clearpriv 一键清空权限
- loopchk 持续检查门的状态`
	//- status 门禁状态显示

	fmt.Println(tip3)
	_, _ = fmt.Scanln(&mode)
	rollback := true
	runonce := false
	for rollback == true {

		rollback = false
		if mode == "open_door" {
			control = []byte{0x9d, 0x10, 0x01, 0x01}
		}
		if mode == "time" {
			locTime := time.Now()
			year := strconv.Itoa(locTime.Year())[2:]
			year2, _ := strconv.Atoi(year)
			control = append(control, uint8(year2))
			control = append(control, uint8(locTime.Month()))
			control = append(control, uint8(locTime.Day()))
			control = append(control, uint8(locTime.Hour()))
			control = append(control, uint8(locTime.Minute()))
			control = append(control, uint8(locTime.Second()))
			break
		}
		if mode == "logcount" {

		}
		if mode == "setdoor" {
			var status uint8
			fmt.Println(tip4)
			_, err := fmt.Scanln(&status)
			if err != nil {
				panic(err)
			}
			control = []byte{0x8f, 0x10, 0x01, status, 0x32, 0x00}
		}
		if mode == "clear" {
			control = []byte{0x8e, 0x10}
		}
		if mode == "readcard" {
			time.Sleep(1000 * time.Millisecond)
			if runonce == false {
				arg2 := int(rec[2])
				arg0 := int(rec[1])*0x100 + int(rec[0])
				fmt.Println("卡号：")
				fmt.Printf("%X\n", arg2*100000+arg0)
				fmt.Println("发卡组织：")
				//fmt.Println([]byte{buff[0],buff[1],buff[2],0})
				fmt.Printf("%X\n", int(rec[2])*0x10000+int(rec[1])*0x100+int(rec[0])*0x1)
				rollback = false
			} else {
				control = []byte{0x8d, 0x10}
				rollback = true
			}
			fmt.Println("状态：")
			fmt.Println(buff[3])
		}
		if mode == "readall" {

		}
		if mode == "addpriv" {
			const tip5 = "请输入卡号（十六进制），如：1ABCDE"
			fmt.Println(tip5)
			var card string
			fmt.Scanln(&card)
			cardid, _ := strconv.ParseUint(card, 16, 32)
			fmt.Println(uint8(math.Mod(math.Mod(float64(cardid), 100000), 0x100)))
			start := convdate(0, 2000, 1, 1, 0, 0, 0)
			end := convdate(0, 2020, 12, 31, 23, 59, 59)
			control = []byte{
				0x07,
				0x11,
				0x01,
				0x00,
				uint8(math.Mod(math.Mod(float64(cardid), 100000), 0x100)),
				uint8(int(math.Mod(float64(cardid), 100000)) / 0x100),
				uint8(cardid / 100000),
				1,
			}
			var buffer bytes.Buffer
			buffer.Write(control)
			buffer.Write(start)
			buffer.Write(end)
			buffer.Write([]byte{1})
			control = buffer.Bytes()
		}
		if mode == "rempriv" {
			const tip5 = "请输入卡号（十六进制），如：1234ABCD"
			fmt.Println(tip5)
			var card string
			fmt.Scanln(&card)
			cardid, _ := strconv.ParseUint(card, 16, 32)
			control = []byte{
				0x08,
				0x11,
				0,
				0,
				uint8(math.Mod(math.Mod(float64(cardid), 100000), 0x100)),
				uint8(int(math.Mod(float64(cardid), 100000)) / 0x100),
				uint8(cardid / 100000),
				1,
			}
		}

		if mode == "status" {

			/*
				//估计没有什么问题，
				//没实践过，鬼知道会返回什么。。。
				//有钱了再把坑给填了，
				//取消注释就能使用了，但是不保证正确！！！

				if runonce == false {
					rollback = true
					control = []byte{0x81, 0x10}
				} else {
					logCount := int(rec[7])*0x100 + int(rec[8])
					permissionsCount := int(rec[10])*0x100 + int(rec[11])
					doorStatus := rec[20]
					warnStatus := rec[22]
					uid := int(rec[14])*100000 + int(rec[12])*0x100 + int(rec[13])
					xingqi := [...]string{"日", "一", "二", "三", "四", "五", "六"}
					fmt.Printf("Log time: 20%02x-%02x-%02x 星期%s %02x:%02x:%02x\n\n", rec[0], rec[1], rec[2], xingqi[rec[3]], rec[4], rec[5], rec[6])
					fmt.Printf("最后一次刷卡卡号%H\n", uid)
					fmt.Println("权限条数：" + strconv.Itoa(permissionsCount))
					fmt.Println("日志条数：" + strconv.Itoa(logCount))
					fmt.Printf("门状态：")
					if doorStatus == 0 {
						fmt.Println("未打开")
					} else {
						fmt.Println("已打开")
					}
					fmt.Println("警告状态：" + string(warnStatus))
					warn := strconv.FormatInt(int64(warnStatus), 2)
					if warn[0] != 0 {
						fmt.Println("遭到胁迫")
					}
					if warn[1] != 0 {
						fmt.Println("门长时间未关闭")
					}
					if warn[2] != 0 {
						fmt.Println("非法闯入")
					}
					if warn[3] != 0 {
						fmt.Println("无效刷卡")
					}

				}
			*/
		}

		if mode == "clearpriv" {
			control = []byte{0x93, 0x10}
		}
		var status2 uint8
		if mode == "loopchk" {
			time.Sleep(1000 * time.Millisecond)
			rollback = true
			control = []byte{0x81, 0x10}
			status := rec[20]
			if status2 != status {
				status2 = status
				if status == 1 {
					fmt.Println("门已打开")
				} else {
					fmt.Println("门已关闭")
				}
			}
		}
		if mode == "resetwarn" {
			control = []byte{0x99, 0x10}
		}
		runonce = true
		buff = prepareData(door, control)
		fmt.Println(buff)
		_, err := udpConn.Write(buff)
		if err != nil {
			fmt.Println("发送失败")
			panic(err)
		} else {
			fmt.Println("发送成功")
		}
	}

}
func prepareData(door []byte, control []byte) []byte {
	data := make([]byte, 34)
	data[0] = 0x7e
	data[33] = 0x0d
	data[1] = door[0]
	data[2] = door[1]
	num := 3
	for _, value := range control {
		data[num] = value
		num++
	}
	data[31], data[32] = checkData(data)
	return data
}
func checkData(data []byte) (byte, byte) {
	checksum := 0
	for i := 1; i <= 30; i++ {
		checksum += int(data[i])
	}
	return uint8(math.Mod(float64(checksum), 0x100)), uint8(checksum / 0x100)
}

func convdate(mode int, year int, month int, day int, hour int, minute int, second int) []byte {
	var ret []byte
	if mode != 0 {
		ret = make([]byte, 4)
		ret[0] = uint8(month / 8 * day)
		ret[1] = uint8(int(math.Mod(float64(year), 100))*2 + month/8)
		ret[2] = uint8(minute/8*32 + second/2)
		ret[3] = uint8(hour*8 + minute/8)
	} else {
		ret = make([]byte, 2)
		ret[0] = uint8(month / 8 * day)
		ret[1] = uint8(int(math.Mod(float64(year), 100))*2 + month/8)
	}
	return ret
}
