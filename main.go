package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	_ "github.com/joho/godotenv/autoload"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

var (
	bought       = false
	purchasing   = false
	lastNotify   = time.Now().Add(-24 * time.Hour)
	chatID       = os.Getenv("CHAT_ID")
	webhook      = os.Getenv("WEBHOOK")
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

func createClient() *common.Client {
	client := sync.OnceValue(func() *common.Client {
		// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
		// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性
		// 以下代码示例仅供参考，建议采用更安全的方式来使用密钥
		// 请参见：https://cloud.tencent.com/document/product/1278/85305
		// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
		credential := common.NewCredential(clientID, clientSecret)
		// 使用临时密钥示例
		// credential := common.NewTokenCredential("SecretId", "SecretKey", "Token")
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
		cpf.HttpProfile.ReqMethod = "POST"
		return common.NewCommonClient(credential, "ap-hongkong", cpf).WithLogger(log.Default())
	})()

	return client
}

func main() {
	if clientID == "" || clientSecret == "" {
		log.Fatal("clientID or clientSecret is empty")
	}
	if chatID == "" || webhook == "" {
		log.Fatal("chat id or webhook env required")
	}
	ctx, cancel := context.WithCancel(context.Background())
	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("main process exit")
			return
		case <-ticker.C:
			queryBundles(ctx)
		}
	}
}

func queryBundles(ctx context.Context) {
	client := createClient()
	request := tchttp.NewCommonRequest("lighthouse", "2020-03-24", "DescribeBundles")
	params := "{\"BundleIds\":[\"bundle_rs_nmc_lin_med1_02\",\"bundle_rs_nmc_lin_med2_01\"]}"

	err := request.SetActionParameters(params)
	if err != nil {
		log.Printf("request set action parameters failed, err:%v\n", err)
		return
	}

	response := tchttp.NewCommonResponse()
	err = client.Send(request, response)
	if err != nil {
		log.Printf("fail to invoke api: %v\n", err)
		return
	}

	var queryBundlesResp DescribeBundlesResp
	if err = sonic.Unmarshal(response.GetBody(), &queryBundlesResp); err != nil {
		log.Printf("fail to unmarshal response: %v\n", err)
		return
	}
	bundles := make(map[string]string, 2)
	for _, bundle := range queryBundlesResp.Response.BundleSet {
		key := fmt.Sprintf("%s-%dC%dG", bundle.BundleTypeDescription, bundle.CPU, bundle.Memory)
		log.Printf("%s: %s\n", key, bundle.BundleSalesState)
		bundles[key] = bundle.BundleSalesState
		if !bought && !purchasing && bundle.BundleSalesState != "SOLD_OUT" {
			createInstance(bundle)
			break
		}
	}
	if time.Since(lastNotify).Hours() >= 1 {
		notify(bundles)
	}
}

func createInstance(bundle Bundle) {
	purchasing = true
	defer func() {
		purchasing = false
	}()
	client := createClient()
	request := tchttp.NewCommonRequest("lighthouse", "2020-03-24", "CreateInstances")
	// bundleID := "bundle_rs_nmc_lin_med1_02"
	rootPassword := os.Getenv("ROOT_PASSWORD")
	if rootPassword == "" {
		rootPassword = "admin@2025"
	}
	params := map[string]any{
		"BundleId":    bundle.BundleID,
		"BlueprintId": "lhbp-mxml4cnq", // Debian 12
		"InstanceChargePrepaid": map[string]any{
			"Period":    1,
			"RenewFlag": "NOTIFY_AND_AUTO_RENEW",
		},
		"InstanceName": "qc6",
		"LoginConfiguration": map[string]any{
			"AutoGeneratePassword": "NO",
			"Password":             rootPassword,
		},
		"AutoVoucher": true,
	}

	err := request.SetActionParameters(params)
	if err != nil {
		log.Printf("fail to set action parameters api: %v\n", err)
		return
	}
	response := tchttp.NewCommonResponse()
	err = client.Send(request, response)
	if err != nil {
		log.Printf("fail to invoke api: %v\n", err)
		return
	}
	var resp CreateInstanceResp
	if err = sonic.Unmarshal(response.GetBody(), &resp); err != nil {
		log.Printf("fail to unmarshal response: %v\n", err)
		return
	}
	if resp.Response.Error != nil {
		log.Printf("fail to invoke api: %v\n", resp.Response.Error)
		return
	}
	bought = true
	bundleName := fmt.Sprintf("%s-%dC%dG", bundle.BundleTypeDescription, bundle.CPU, bundle.Memory)
	notifyBought(bundleName)
	log.Printf("CreateInstanceResp: %+v\n", resp.Response.InstanceIDSet)
}

func notify(bundles map[string]string) {
	markdownContent := "## ⚙️ **监控服务运行中**\n"
	for k, v := range bundles {
		markdownContent += fmt.Sprintf("- **%s**: %s\n", k, v)
	}
	markdownContent += "\n\n**通知时间**：" + time.Now().Format(time.DateTime)
	payload := map[string]any{
		"chatid":  chatID,
		"msgtype": "markdown",
		"markdown": map[string]any{
			"content": markdownContent,
		},
	}
	body, _ := sonic.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("fail to invoke api: %v\n", err)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode == http.StatusOK {
		lastNotify = time.Now()
	}
	fmt.Println("response Status:", resp.Status)
}

func notifyBought(bundle string) {
	markdownContent := "## ✅ **锐驰自动购买成功**\n"
	markdownContent += fmt.Sprintf("- **型号**: %s\n", bundle)
	markdownContent += "\n\n**通知时间**：" + time.Now().Format(time.DateTime)
	payload := map[string]any{
		"chatid":  chatID,
		"msgtype": "markdown",
		"markdown": map[string]any{
			"content": markdownContent,
		},
	}
	body, _ := sonic.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("fail to invoke api: %v", err)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	log.Printf("response Status: %v", resp.Status)
}
