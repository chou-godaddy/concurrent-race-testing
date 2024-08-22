package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	msmqclient "github.com/gdcorp-domains/fulfillment-generic-queue-client/client"
	registrymodels "github.com/gdcorp-domains/fulfillment-goapimodels/registry"
	httpclient "github.com/gdcorp-domains/fulfillment-golang-httpclient"
	logging "github.com/gdcorp-domains/fulfillment-golang-logging"
	sso "github.com/gdcorp-domains/fulfillment-golang-sso-auth/client"
	"github.com/gdcorp-domains/fulfillment-golang-sso-auth/modifiers"
	"github.com/gdcorp-domains/fulfillment-golang-sso-auth/token/expiry"
	registrydomains "github.com/gdcorp-domains/fulfillment-registry-domains"
	registrydomainsmodel "github.com/gdcorp-domains/fulfillment-registry-domains/models"
	golanghttpclient "github.com/gdcorp-domains/fulfillment-golang-httpclient"
	gosecrets "github.com/gdcorp-domains/fulfillment-gosecrets"
	dnsapiclient "github.com/gdcorp-domains/fulfillment-golang-clients/dnsapi"

	"log"
	"net/url"
	"os"
	"runtime/trace"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Touchpoint struct {
	ResourceID   int
	DomainName   string
	ResourceType string
	Message      string
	Time         string
	ServerName   string
	ServiceName  string
	LogFile      string
}

var goroutineCounter int64

func generateGoroutineID() int64 {
	return atomic.AddInt64(&goroutineCounter, 1)
}

type TaskFunc func() error

func RunConcurrently(tasks []TaskFunc) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(tsk TaskFunc, id int64) {
			defer wg.Done()
			start := time.Now()
			fmt.Printf("Goroutine ID: %d started\n", id)
			if err := tsk(); err != nil {
				errCh <- err
			}
			fmt.Printf("Goroutine ID: %d completed in %v\n", id, time.Since(start))
		}(task, generateGoroutineID())
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return <-errCh
	}
	return nil
}

func ProcessTouchpoints(ctx context.Context, touchpoints []Touchpoint, ruleName string, msmqClient msmqclient.MSMQClient) error {
	start := time.Now()
	defer func() {
		log.Printf("ProcessTouchpoints completed in %v\n", time.Since(start))
	}()
	tasks := make([]TaskFunc, len(touchpoints))
	for i, touchpoint := range touchpoints {
		tp := touchpoint // Create a local copy to avoid closure issues
		tasks[i] = func() error {
			return processSingleTouchpoint(ctx, msmqClient, ruleName, tp)
		}
	}

	return RunConcurrently(tasks)
}

func processSingleTouchpoint(ctx context.Context, msmqClient msmqclient.MSMQClient, ruleName string, touchpoint Touchpoint) error {
	fmt.Printf("Processing touchpoint id %s\n", touchpoint.Message)
	b, err := json.Marshal(touchpoint)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal touchpoint with error, %v", ruleName, err)
	}
	msmqLabel := fmt.Sprintf("%s-%d", touchpoint.ServiceName, touchpoint.ResourceID)
	msmqRes, err := msmqClient.SendMessage(ctx, "regloggingsvc-incoming", msmqLabel, string(b))
	if msmqRes == nil || msmqRes.Body.SendMessageResponse == nil {
		errStr := "No response received from GenericQueueWebSvc"
		if err != nil {
			errStr = fmt.Sprintf("err: %s", err.Error())
		}
		return fmt.Errorf("%s: failed to send touchpoint with error, %v", ruleName, errStr)
	}
	msmqResBody := string(msmqRes.Body.SendMessageResponse.SendMessageResult.Text)
	if !strings.Contains(msmqResBody, "success") {
		return fmt.Errorf("%s: failed to send touchpoint for resourceId %d with error, %s", ruleName, touchpoint.ResourceID, msmqResBody)
	}
	return nil
}

func ProcessDomainInfoCalls(ctx context.Context, domainNames []string, registryDomainsClient registrydomains.Client) error {
	start := time.Now()
	defer func() {
		log.Printf("ProcessDomainInfoCalls completed in %v\n", time.Since(start))
	}()
	tasks := make([]TaskFunc, len(domainNames))
	for i, name := range domainNames {
		n := name
		tasks[i] = func() error {
			return processSingleDomainInfoCall(ctx, registryDomainsClient, n)
		}
	}

	return RunConcurrently(tasks)
}

func ProcessDomainInfoCallsSequentially(ctx context.Context, domainNames []string, registryDomainsClient registrydomains.Client) error {
	start := time.Now()
	defer func() {
		log.Printf("ProcessDomainInfoCallsSequentially completed in %v\n", time.Since(start))
	}()
	for _, name := range domainNames {
		if err := processSingleDomainInfoCall(ctx, registryDomainsClient, name); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func processSingleDomainInfoCall(ctx context.Context, registryDomainsClient registrydomains.Client, domainName string) error {
	fmt.Printf("Processing domain info for domain %s\n", domainName)
	infoReq := registrydomainsmodel.RegistryDomainInfo{
		Keys: registrymodels.Keys{
			RGRegistrar: strings.ToLower("GODADDY"),
			TLD:         strings.ToLower("org"),
		},
		HostsAttr: registrydomainsmodel.AllHosts,
		Domain:    domainName,
		Fields:    []registrydomainsmodel.DomainInfoField{registrydomainsmodel.Status, registrydomainsmodel.NameServer},
	}
	resp, _, err := registryDomainsClient.GetRegistryDomainsInfo(context.Background(), infoReq, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	return nil
}

func ProcessDNSGetDefaultNameservers(ctx context.Context, domainNames []string, dnsAPIClient dnsapiclient.Client) error {
	start := time.Now()
	defer func() {
		log.Printf("ProcessDNSGetDefaultNameservers completed in %v\n", time.Since(start))
	}()
	tasks := make([]TaskFunc, len(domainNames))
	for i, name := range domainNames {
		n := name
		tasks[i] = func() error {
			return processSingleDNSGetDefaultNameservers(ctx, dnsAPIClient, n)
		}
	}

	return RunConcurrently(tasks)
}

func processSingleDNSGetDefaultNameservers(ctx context.Context, dnsAPIClient dnsapiclient.Client, domainName string) error {
	fmt.Printf("Processing default nameservers for domain %s\n", domainName)
	resp, code, err := dnsAPIClient.GetDomainDefaultNameservers(context.Background(), "32728873", "1", domainName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Response: %v, Code: %d\n", resp, code)
	return nil
}

var counter int

func increment(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		counter++
	}
}

func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}

func main() {
	// Add to get the race report purposely
	// start := time.Now()
	// var t *time.Timer
	// t = time.AfterFunc(randomDuration(), func() {
	// 	fmt.Println(time.Now().Sub(start))
	// 	t.Reset(randomDuration())
	// })
	// time.Sleep(5 * time.Second)


	// HttpClient Plain
	httpconfig := &httpclient.ClientConfig{
		MaxConnections: 10,
		TimeoutMs:      90000,
	}
	httpClient, err := httpclient.NewClientBuilder(nil).HTTPClient(context.Background(), httpconfig)
	if err != nil {
		panic(err)
	}

	// HttpClient with IAM
	iamRefresher := sso.NewIAMRefresher(sso.NewIAMClient(sso.IAMClientConfig{SSOURL: "sso-host", PrimaryRegion: "us-west-2", SecondaryRegions: []string{"us-west-2"}}, httpClient), 15*time.Minute, logging.New(logging.Config{}))
	httpClientWithIAM := httpClient.WithRequestModifier(modifiers.NewAddIAMJWTModifier(expiry.Medium, iamRefresher))

	var certLoader golanghttpclient.CertGetter

	secretRetriever := gosecrets.NewSecretRetriever()
	certLoader = golanghttpclient.NewCertGetter(
		golanghttpclient.CertConfig{
			Cert: gosecrets.SecretConfig{
				AWS: &gosecrets.AWSSecretConfig{
					Name:   "cert",
					Region: "us-west-2",
				},
			},
			Key: gosecrets.SecretConfig{
				AWS: &gosecrets.AWSSecretConfig{
					Name:   "key",
					Region: "us-west-2",
				},
			},
			Chain: gosecrets.SecretConfig{
				AWS: &gosecrets.AWSSecretConfig{
					Name:   "chain",
					Region: "us-west-2",
				},
			},
			Root: gosecrets.SecretConfig{
				AWS: &gosecrets.AWSSecretConfig{
					Name:   "ca-bundle.crt",
					Region: "us-west-2",
				},
			},
		}, secretRetriever)

	clientConfig := &golanghttpclient.ClientConfig{
		MaxConnections: 100,
		TimeoutMs:      10000,
	}

	httpClientBuilder := golanghttpclient.NewClientBuilder(certLoader)
	httpClientWithCert, err := httpClientBuilder.HTTPClient(context.Background(), clientConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Example of using httpclient plain
	msmqurl := "msmq-test-url"
	msmqQueueURL, _ := url.Parse(msmqurl)
	MSMQClient := msmqclient.NewMSMQClient(msmqQueueURL, httpClient)
	touchpoints := []Touchpoint{}
	for i := 0; i < 20; i++ {
		touchpoints = append(touchpoints, Touchpoint{
			ResourceID:   32728958,
			ResourceType: "domain",
			DomainName:   "example.com",
			Message:      fmt.Sprintf("id: %d", i),
			Time:         time.Date(2023, 07, 07, 0, 0, 0, 0, time.UTC).String(),
		})
	}
	err = ProcessTouchpoints(context.Background(), touchpoints, "NameserverUpdate", MSMQClient)
	if err != nil {
		fmt.Println(err)
	}

	// Example of using httpclient with IAM
	registryDomainsurl := "domains-api-test-url"
	registryDomainsURL, _ := url.Parse(registryDomainsurl)
	RegistryDomainsClient := registrydomains.NewClient(registryDomainsURL, httpClientWithIAM)

	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace file: %v", err)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	var domainNames []string
	for i := 100; i < 120; i++ {
		domainNames = append(domainNames, fmt.Sprintf("authcodegen%d.org", i))
	}

	err = ProcessDomainInfoCalls(context.Background(), domainNames, RegistryDomainsClient)

	// Example of using httpclient with Cert
	dnsAPIURLClientURL, _ := url.Parse("https://dns.api.int.test-godaddy.com")
	DNSAPIClient := dnsapiclient.NewClient(dnsAPIURLClientURL, httpClientWithCert)
	err = ProcessDNSGetDefaultNameservers(context.Background(), domainNames, DNSAPIClient)
	if err != nil {
		fmt.Println(err)
	}
}
