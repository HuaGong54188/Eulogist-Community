package Eulogist

import (
	"Eulogist/core/tools/skin_process"
	Client "Eulogist/proxy/mc_client"
	Server "Eulogist/proxy/mc_server"
	"Eulogist/proxy/persistence_data"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"sync"
	"time"

	"github.com/pterm/pterm"
)

// Eulogist 函数是整个“赞颂者”程序的入口点
func Eulogist() error {
	var err error
	var config *EulogistConfig
	var neteaseConfigPath string
	var waitGroup sync.WaitGroup
	var client *Client.MinecraftClient
	var server *Server.MinecraftServer
	var clientWasConnected chan struct{}
	var persistenceData *persistence_data.PersistenceData = new(persistence_data.PersistenceData)

	// 读取配置文件
	{
	pterm.DefaultBox.Println(pterm.LightCyan("https://github.com/HuaGong54188/Eulogist-Community/tree/main"))
	pterm.Println(pterm.Yellow("注 - 这是 Eulogist 的非官方汉化版本"))
		config, err = ReadEulogistConfig()
		if err != nil {
			return fmt.Errorf("Eulogist: %v", err)
		}
	}

	// 使赞颂者连接到网易租赁服
	{
		pterm.Info.Println("现在，我们来尝试连接到验证服务器。")

		server, err = Server.ConnectToServer(
			Server.BasicConfig{
				ServerCode:     config.RentalServerCode,
				ServerPassword: config.RentalServerPassword,
				Token:          config.FBToken,
				AuthServer:     LookUpAuthServerAddress(config.FBToken),
			},
			persistenceData,
		)
		if err != nil {
			return fmt.Errorf("赞颂者: %v", err)
		}
		defer server.Conn.CloseConnection()

		pterm.Success.Println("喜报！成功连接到网易租赁服，接下来我们尝试与租赁服握手。")

		err = server.FinishHandshake()
		if err != nil {
			return fmt.Errorf("赞颂者: %v", err)
		}

		pterm.Success.Println("与租赁服握手成功了。")
	}

	// 根据配置文件的启动类型决定启动方式
	if config.LaunchType == LaunchTypeNormal {
		// 初始化
		var playerSkin *skin_process.Skin
		var neteaseSkinFileName string
		var skinIsSlim bool
		var useAccountSkin bool
		// 检查 Minecraft 客户端是否存在
		if !FileExist(config.NEMCPath) {
			return fmt.Errorf("赞颂者: 找不到我的世界客户端，可能你没有下载或者路径错误")
		}
		// 取得皮肤数据
		playerSkin = server.PersistenceData.SkinData.NeteaseSkin
		useAccountSkin = (!FileExist(config.SkinPath) && playerSkin != nil)
		// 皮肤处理
		if useAccountSkin {
			// 生成皮肤文件
			if skin_process.IsZIPFile(playerSkin.FullSkinData) {
				neteaseSkinFileName = "skin.zip"
			} else {
				neteaseSkinFileName = "skin.png"
			}
			err = os.WriteFile(neteaseSkinFileName, playerSkin.FullSkinData, 0600)
			if err != nil {
				return fmt.Errorf("赞颂者: %v", err)
			}
			currentPath, _ := os.Getwd()
			config.SkinPath = fmt.Sprintf(`%s\%s`, currentPath, neteaseSkinFileName)
			// 皮肤纤细处理
			skinIsSlim = playerSkin.SkinIsSlim
		}
		// 启动 Eulogist 服务器
		client, clientWasConnected, err = Client.RunServer(persistenceData)
		if err != nil {
			return fmt.Errorf("赞颂者: %v", err)
		}
		defer client.Conn.CloseConnection()
		// 生成网易配置文件
		neteaseConfigPath, err = GenerateNetEaseConfig(config.SkinPath, skinIsSlim, client.Address.IP.String(), client.Address.Port)
		if err != nil {
			return fmt.Errorf("赞颂者: %v", err)
		}
		// 启动 Minecraft 客户端
		command := exec.Command(config.NEMCPath, fmt.Sprintf("config=%s", neteaseConfigPath))
		go command.Run()
		// 打印准备完成的信息
		pterm.Success.Println("赞颂者已经迫不及待了！接下来会自动启动我的世界客户端，你无需手动进入。\n随后，客户端将自动连接到赞颂者。")
	} else {
		// 启动 Eulogist 服务器
		client, clientWasConnected, err = Client.RunServer(persistenceData)
		if err != nil {
			return fmt.Errorf("Eulogist: %v", err)
		}
		defer client.Conn.CloseConnection()
		// 打印赞颂者准备完成的信息
		pterm.Success.Printf(
			"赞颂者已准备完成，接下来需要你手动连接\n连接到赞颂者的地址(当成正常的服务器IP连接即可): %s:%d\n",
			client.Address.IP.String(), client.Address.Port,
		)
	}

	// 等待 Minecraft 客户端与赞颂者完成基本数据包交换
	{
		// 等待 Minecraft 客户端连接
		if config.LaunchType == LaunchTypeNormal {
			timer := time.NewTimer(time.Second * 120)
			defer timer.Stop()
			select {
			case <-timer.C:
				return fmt.Errorf("赞颂者: 无法与Minecraft客户端建立连接")
			case <-clientWasConnected:
				close(clientWasConnected)
			}
		} else {
			<-clientWasConnected
			close(clientWasConnected)
		}
		pterm.Success.Println("与客户端连接成功！现在我们尝试与其握手。")
		// 等待 Minecraft 客户端完成握手
		err = client.WaitClientHandshakeDown()
		if err != nil {
			return fmt.Errorf("Eulogist: %v", err)
		}
		pterm.Success.Println("与客户端握手成功后，你就能登录到网易租赁服了。")
	}

	// 设置等待队列
	waitGroup.Add(2)

	// 处理网易租赁服到赞颂者的数据包
	go func() {
		// 关闭已建立的所有连接
		defer func() {
			waitGroup.Add(-1)
			server.Conn.CloseConnection()
			client.Conn.CloseConnection()
		}()
		// 显示程序崩溃错误信息
		defer func() {
			r := recover()
			if r != nil {
				pterm.Error.Printf(
					"Eulogist/GoFunc/RentalServerToEulogist: err = %v\n\n[Stack Info]\n%s\n",
					r, string(debug.Stack()),
				)
				fmt.Println()
			}
		}()
		// 数据包抄送
		for {
			// 初始化一个函数，
			// 用于同步数据到 Minecraft 客户端
			syncFunc := func() error {
				if shieldID := server.Conn.GetShieldID().Load(); shieldID != 0 {
					client.Conn.GetShieldID().Store(shieldID)
				}
				return nil
			}
			// 读取、过滤数据包，
			// 然后抄送其到 Minecraft 客户端
			errResults, syncError := server.FiltePacketsAndSendCopy(server.Conn.ReadPackets(), client.Conn.WritePackets, syncFunc)
			if syncError != nil {
				pterm.Warning.Printf("赞颂者: 处理数据包的时候数据同步失败了，日志为 %v", syncError)
			}
			for _, err = range errResults {
				if err != nil {
					pterm.Warning.Printf("赞颂者: 处理来自服务器时错误的数据包: %v\n", err)
				}
			}
			// 检查连接状态
			select {
			case <-server.Conn.GetContext().Done():
				return
			case <-client.Conn.GetContext().Done():
				return
			default:
			}
		}
	}()

	// 处理 Minecraft 客户端到赞颂者的数据包
	go func() {
		// 关闭已建立的所有连接
		defer func() {
			waitGroup.Add(-1)
			client.Conn.CloseConnection()
			server.Conn.CloseConnection()
		}()
		// 显示程序崩溃错误信息
		defer func() {
			r := recover()
			if r != nil {
				pterm.Error.Printf(
					"Eulogist/GoFunc/MinecraftClientToEulogist: err = %v\n\n[Stack Info]\n%s\n",
					r, string(debug.Stack()),
				)
				fmt.Println()
			}
		}()
		// 数据包抄送
		for {
			// 初始化一个函数，
			// 用于同步数据到网易租赁服
			syncFunc := func() error {
				return nil
			}
			// 读取、过滤数据包，
			// 然后抄送其到网易租赁服
			errResults, syncError := client.FiltePacketsAndSendCopy(client.Conn.ReadPackets(), server.Conn.WritePackets, syncFunc)
			if syncError != nil {
				pterm.Warning.Printf("赞颂者: 处理来自客户端的数据时，同步数据失败了，日志为 %v", syncError)
			}
			for _, err = range errResults {
				if err != nil {
					pterm.Warning.Printf("赞颂者: 处理来自客户端时错误的数据包: %v\n", err)
				}
			}
			// 检查连接状态
			select {
			case <-client.Conn.GetContext().Done():
				return
			case <-server.Conn.GetContext().Done():
				return
			default:
			}
		}
	}()

	// 等待所有 goroutine 完成
	waitGroup.Wait()
	pterm.Info.Println("现在关闭了所有的链接，如果还想再次进服，请重复一次之前的步骤。")
	return nil
}
