package sys

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	INTERNET_PER_CONN_FLAGS               = 1
	INTERNET_PER_CONN_PROXY_SERVER        = 2
	INTERNET_PER_CONN_PROXY_BYPASS        = 3
	INTERNET_OPTION_REFRESH               = 37
	INTERNET_OPTION_SETTINGS_CHANGED      = 39
	INTERNET_OPTION_PER_CONNECTION_OPTION = 75
)

type INTERNET_PER_CONN_OPTION struct {
	dwOption uint32
	dwValue  uint64 // 注意 32位 和 64位 struct 和 union 内存对齐
}

type INTERNET_PER_CONN_OPTION_LIST struct {
	dwSize        uint32
	pszConnection *uint16
	dwOptionCount uint32
	dwOptionError uint32
	pOptions      uintptr
}

func set(proxy string) error {
	winInet, err := syscall.LoadLibrary("Wininet.dll")
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("LoadLibrary Wininet.dll Error: %s", err))
	}
	InternetSetOptionW, err := syscall.GetProcAddress(winInet, "InternetSetOptionW")
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("GetProcAddress InternetQueryOptionW Error: %s", err))
	}

	options := [3]INTERNET_PER_CONN_OPTION{}
	options[0].dwOption = INTERNET_PER_CONN_FLAGS
	if proxy == "" {
		options[0].dwValue = 1
	} else {
		options[0].dwValue = 2
	}
	options[1].dwOption = INTERNET_PER_CONN_PROXY_SERVER

	px, err := syscall.UTF16PtrFromString(proxy)
	options[1].dwValue = uint64(uintptr(unsafe.Pointer(px)))
	options[2].dwOption = INTERNET_PER_CONN_PROXY_BYPASS
	if err != nil {
		return err
	}

	hstr, err := syscall.UTF16PtrFromString("localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;172.32.*;192.168.*")

	if err != nil {
		return err
	}

	options[2].dwValue = uint64(uintptr(unsafe.Pointer(hstr)))

	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.pszConnection = nil
	list.dwOptionCount = 3
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))

	// https://www.cnpython.com/qa/361707
	callInternetOptionW := func(dwOption uintptr, lpBuffer uintptr, dwBufferLength uintptr) error {
		r1, _, err := syscall.Syscall6(InternetSetOptionW, 4, 0, dwOption, lpBuffer, dwBufferLength, 0, 0)
		if r1 != 1 {
			return err
		}
		return nil
	}

	err = callInternetOptionW(INTERNET_OPTION_PER_CONNECTION_OPTION, uintptr(unsafe.Pointer(&list)), uintptr(unsafe.Sizeof(list)))
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_PER_CONNECTION_OPTION Error: %s", err)
	}
	err = callInternetOptionW(INTERNET_OPTION_SETTINGS_CHANGED, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_SETTINGS_CHANGED Error: %s", err)
	}
	err = callInternetOptionW(INTERNET_OPTION_REFRESH, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_REFRESH Error: %s", err)
	}
	return nil
}

func setup(b bool) {
	proxy := "127.0.0.1:1080"
	if b {
		if err := set(proxy); err == nil {
			fmt.Println("开启代理成功!")
		} else {
			fmt.Println("开启代理失败 (+﹏+) ")
		}
	} else {
		if err := set(""); err == nil {
			fmt.Println("关闭代理成功!")
		} else {
			fmt.Println("关闭代理失败 (+﹏+)")
		}

	}
}

func SetProxy(b bool) {

	if b {
		setup(b)
	}

	for {
		var ok string
		fmt.Print("1开启  2关闭  3关闭并退出 =>")
		//当程序只是到fmt.Scanln(&name)程序会停止执行等待用户输入
		fmt.Scanln(&ok)

		index, err := strconv.Atoi(ok)

		if err != nil {
			continue
		}

		if index == 1 {
			setup(true)
		} else if index == 2 {
			setup(false)
		} else if index == 3 {
			os.Exit(0)
		}

	}

}
