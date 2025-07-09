package app

import (
	"errors"
	"go-svc-tpl/app/controller"
	"go-svc-tpl/logs"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"github.com/LagrangeDev/LagrangeGo/utils"
)

var qqclient *client.QQClient

func GetQQClient() *client.QQClient {
	return qqclient
}

func InitLagrangeBot() {
	// 使用特定的协议版本
	appInfo := auth.AppList["linux"]["3.2.15-30366"]
	// 创建设备信息
	deviceInfo := &auth.DeviceInfo{
		GUID:          "cfcd208495d565ef66e7dff9f98764da",
		DeviceName:    "Lagrange-DCFCD07E",
		SystemKernel:  "Windows 10.0.22631",
		KernelVersion: "10.0.22631",
	}

	// 创建qqclient实例
	qqclient = client.NewClient(0, "")
	// 使用协议版本
	qqclient.UseVersion(appInfo)
	// 添加signserver，注意要和appinfo版本匹配
	qqclient.AddSignServer("https://sign.lagrangecore.org/api/sign/30366")
	// 使用设备信息
	qqclient.UseDevice(deviceInfo)

	// 从保存的sig.bin文件读取登录信息
	data, err := os.ReadFile("sig.bin")
	if err != nil {
		logrus.Warn("read sig error:", err)
	} else {
		// 将登录信息反序列化
		sig, err := auth.UnmarshalSigInfo(data, true)
		if err != nil {
			logrus.Warn("load sig error:", err)
		} else {
			// 如果登录信息有效，则使用登录信息登录
			qqclient.UseSig(sig)
		}
	}

	controller.ClassifyReq(qqclient)
}

func StartLagrangeBot() {
	err := func(c *client.QQClient, passwordLogin bool) error {
		logs.Info("login with password")
		// 如果登录信息存在，可以使用fastlogin
		err := c.FastLogin()
		if err == nil {
			return nil
		}

		if passwordLogin {
			// 密码登录，目前无法使用
			ret, err := c.PasswordLogin()
			for {
				if err != nil {
					logs.Error("密码登录失败: %s", err)
					break
				}
				if ret.Success {
					return nil
				}
				switch ret.Error {
				case client.SliderNeededError:
					logs.Warn("captcha verification required")
					logs.Warn(ret.VerifyURL)
					aid := strings.Split(strings.Split(ret.VerifyURL, "sid=")[1], "&")[0]
					logs.Warn("ticket?->")
					ticket := utils.ReadLine()
					logs.Warn("rand_str?->")
					randStr := utils.ReadLine()
					ret, err = c.SubmitCaptcha(ticket, randStr, aid)
					continue
				case client.UnsafeDeviceError:
					vf, err := c.GetNewDeviceVerifyURL()
					if err != nil {
						return err
					}
					logs.Info(vf)
					err = c.NewDeviceVerify(vf)
					if err != nil {
						return err
					}
				default:
					logs.Error("Unhandled exception raised: %s", ret.ErrorMessage)
				}
			}
		}
		logs.Info("login with qrcode")

		// 扫码登录流程
		// 首先获取二维码
		png, _, err := c.FetchQRCodeDefault()
		if err != nil {
			return err
		}
		qrcodePath := "qrcode.png"
		// 保存到本地以供扫码
		err = os.WriteFile(qrcodePath, png, 0666)
		if err != nil {
			return err
		}
		logs.Info("qrcode saved to %s", qrcodePath)
		for {
			// 轮询二维码扫描结果
			retCode, err := c.GetQRCodeResult()
			if err != nil {
				logs.Error(err.Error())
				return err
			}
			// 等待扫码
			if retCode.Waitable() {
				time.Sleep(3 * time.Second)
				continue
			}
			if !retCode.Success() {
				return errors.New(retCode.Name())
			}
			break
		}
		// 扫码完成后就可以进行登录
		_, err = c.QRCodeLogin()
		return err
	}(qqclient, false)

	if err != nil {
		logs.Error("login err:", err)
		return
	}
	logs.Info("login successed")

	defer qqclient.Release()

	defer func() {
		// 序列化登录信息以便下次使用
		data, err := qqclient.Sig().Marshal()
		if err != nil {
			logs.Error("marshal sig.bin err:", err)
			return
		}
		err = os.WriteFile("sig.bin", data, 0644)
		if err != nil {
			logs.Error("write sig.bin err:", err)
			return
		}
		logs.Info("sig saved into sig.bin")
	}()

	// setup the main stop channel
	mc := make(chan os.Signal, 2)
	signal.Notify(mc, os.Interrupt, syscall.SIGTERM)
	for {
		switch <-mc {
		case os.Interrupt, syscall.SIGTERM:
			return
		}
	}
}
