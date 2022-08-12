package serviceDiscovery

import (
	"fmt"
	"sync"
	"time"
)

var (
	zkUrl = []string{"192.168.31.202:2181"}
)

func handlerNodeEvent(endpointId string, server ServiceDiscovery) {
	for {
		select {
		case events, ok := <-server.NodeEvent():
			if ok {
				for _, event := range events {
					if event.Et == ETAdd {
						fmt.Printf("endpoint add.{cEndpointId : %s} {id : %s}, {content : %s}\n", endpointId, event.Endpoint.Id(), event.Endpoint.Content())
					} else if event.Et == ETReconnect {
						fmt.Printf("endpoint reconnect.{cEndpointId : %s}  {id : %s}, {content : %s}\n", endpointId, event.Endpoint.Id(), event.Endpoint.Content())
					} else if event.Et == ETDisconnect {
						fmt.Printf("endpoint disconnect.{cEndpointId : %s}  {id : %s}, {content : %s}\n", endpointId, event.Endpoint.Id(), event.Endpoint.Content())
					}
				}
			} else {
				fmt.Printf("Service discovery close. {cEndpointId : %s}\n", endpointId)
				return
			}
		}
	}
}

// ExampleNormalStartupNode
func ExampleNormalStartupNode() {
	endpointData := []struct {
		id       string
		waitTime time.Duration
	}{
		{
			id:       "1.1.2.1",
			waitTime: time.Second * 1,
		},
		{
			id:       "1.1.2.2",
			waitTime: time.Second * 2,
		},
		{
			id:       "1.1.3.1",
			waitTime: time.Second * 3,
		},
	}
	wg := sync.WaitGroup{}
	for _, endpoint := range endpointData {
		wg.Add(1)
		testdata := endpoint
		go func() {
			defer wg.Done()
			sd, err := NewServiceDiscovery(
				zkUrl,
				testdata.id,
				WithOpenDetailInfo(false),
				WithEndpointEventNotify(true),
			)
			if err != nil {
				fmt.Println(err)
			}
			go handlerNodeEvent(testdata.id, sd)
			time.Sleep(testdata.waitTime)
			sd.Close()
		}()
	}
	wg.Wait()
	//Unordered output:
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.2}, {content : }
	//Service discovery close. {cEndpointId : 1.1.2.1}
	//endpoint disconnect.{cEndpointId : 1.1.3.1}  {id : 1.1.2.1}, {content : }
	//endpoint disconnect.{cEndpointId : 1.1.2.2}  {id : 1.1.2.1}, {content : }
	//Service discovery close. {cEndpointId : 1.1.2.2}
	//endpoint disconnect.{cEndpointId : 1.1.3.1}  {id : 1.1.2.2}, {content : }

}

func ExampleReconnectEndpoint() {
	endpointData := []struct {
		id       string
		waitTime time.Duration
	}{
		{
			id:       "1.1.2.1",
			waitTime: time.Second * 3,
		},
		{
			id:       "1.1.2.2",
			waitTime: time.Second * 2,
		},
		{
			id:       "1.1.3.1",
			waitTime: time.Second,
		},
	}
	wg := sync.WaitGroup{}
	for _, data := range endpointData {
		wg.Add(1)
		testData := data
		go func() {
			defer wg.Done()
			sd, err := NewServiceDiscovery(
				zkUrl,
				testData.id,
				WithOpenDetailInfo(false),
				WithEndpointEventNotify(true),
			)
			if err != nil {
				fmt.Println(err)
			}
			go handlerNodeEvent(testData.id, sd)
			time.Sleep(testData.waitTime)
			sd.Close()
			// 重新启动
			sd, err = NewServiceDiscovery(
				zkUrl,
				testData.id,
				WithOpenDetailInfo(false),
				WithEndpointEventNotify(true),
			)
			if err != nil {
				fmt.Println(err)
			}
			go handlerNodeEvent(testData.id, sd)
			time.Sleep(testData.waitTime)
			sd.Close()
		}()
	}
	wg.Wait()
	//Unordered output:
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.1}, {content : }
	//Service discovery close. {cEndpointId : 1.1.3.1}
	//endpoint disconnect.{cEndpointId : 1.1.2.2}  {id : 1.1.3.1}, {content : }
	//endpoint disconnect.{cEndpointId : 1.1.2.1}  {id : 1.1.3.1}, {content : }
	//endpoint reconnect.{cEndpointId : 1.1.2.2}  {id : 1.1.3.1}, {content : }
	//endpoint reconnect.{cEndpointId : 1.1.2.1}  {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.1}, {content : }
	//Service discovery close. {cEndpointId : 1.1.2.2}
	//endpoint disconnect.{cEndpointId : 1.1.3.1}  {id : 1.1.2.2}, {content : }
	//endpoint disconnect.{cEndpointId : 1.1.2.1}  {id : 1.1.2.2}, {content : }
	//Service discovery close. {cEndpointId : 1.1.3.1}
	//endpoint disconnect.{cEndpointId : 1.1.2.1}  {id : 1.1.3.1}, {content : }
	//endpoint reconnect.{cEndpointId : 1.1.2.1}  {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.1}, {content : }
	//Service discovery close. {cEndpointId : 1.1.2.1}
	//endpoint disconnect.{cEndpointId : 1.1.2.2}  {id : 1.1.2.1}, {content : }
	//endpoint reconnect.{cEndpointId : 1.1.2.2}  {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.1}, {content : }
	//Service discovery close. {cEndpointId : 1.1.2.2}
	//endpoint disconnect.{cEndpointId : 1.1.2.1}  {id : 1.1.2.2}, {content : }

}

// ExampleWithRootPath
func ExampleWithRootPath() {
	rootPath := "/customRootPath"
	endpointData := []struct {
		id       string
		waitTime time.Duration
	}{
		{
			id:       "1.1.2.1",
			waitTime: time.Second * 1,
		},
		{
			id:       "1.1.2.2",
			waitTime: time.Second * 2,
		},
		{
			id:       "1.1.3.1",
			waitTime: time.Second * 3,
		},
	}
	wg := sync.WaitGroup{}
	for _, endpoint := range endpointData {
		wg.Add(1)
		testdata := endpoint
		go func() {
			defer wg.Done()
			sd, err := NewServiceDiscovery(
				zkUrl,
				testdata.id,
				WithOpenDetailInfo(false),
				WithEndpointEventNotify(true),
				WithRootPath(rootPath),
			)
			if err != nil {
				fmt.Println(err)
			}
			go handlerNodeEvent(testdata.id, sd)
			time.Sleep(testdata.waitTime)
			sd.Close()
		}()
	}
	wg.Wait()
	//Unordered output:
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.2}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.3.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.1}, {content : }
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.2}, {content : }
	//Service discovery close. {cEndpointId : 1.1.2.1}
	//endpoint disconnect.{cEndpointId : 1.1.3.1}  {id : 1.1.2.1}, {content : }
	//endpoint disconnect.{cEndpointId : 1.1.2.2}  {id : 1.1.2.1}, {content : }
	//Service discovery close. {cEndpointId : 1.1.2.2}
	//endpoint disconnect.{cEndpointId : 1.1.3.1}  {id : 1.1.2.2}, {content : }

}

func ExampleWithNodeContent() {
	rootPath := "/customRootPath"
	rootContent := "hello world"
	endpointData := []struct {
		id       string
		waitTime time.Duration
	}{
		{
			id:       "1.1.2.1",
			waitTime: time.Second * 1,
		},
		{
			id:       "1.1.2.2",
			waitTime: time.Second * 2,
		},
		{
			id:       "1.1.3.1",
			waitTime: time.Second * 3,
		},
	}
	wg := sync.WaitGroup{}
	for _, endpoint := range endpointData {
		wg.Add(1)
		testdata := endpoint
		go func() {
			defer wg.Done()
			sd, err := NewServiceDiscovery(
				zkUrl,
				testdata.id,
				WithOpenDetailInfo(false),
				WithEndpointEventNotify(true),
				WithRootPath(rootPath),
				WithNodeContent(rootContent),
			)
			if err != nil {
				fmt.Println(err)
			}
			go handlerNodeEvent(testdata.id, sd)
			time.Sleep(testdata.waitTime)
			sd.Close()
		}()
	}
	wg.Wait()
	//Unordered output:
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.3.1}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.1}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.2.1} {id : 1.1.2.2}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.3.1}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.1}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.2.2} {id : 1.1.2.2}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.3.1}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.1}, {content : hello world}
	//endpoint add.{cEndpointId : 1.1.3.1} {id : 1.1.2.2}, {content : hello world}
	//Service discovery close. {cEndpointId : 1.1.2.1}
	//endpoint disconnect.{cEndpointId : 1.1.3.1}  {id : 1.1.2.1}, {content : hello world}
	//endpoint disconnect.{cEndpointId : 1.1.2.2}  {id : 1.1.2.1}, {content : hello world}
	//Service discovery close. {cEndpointId : 1.1.2.2}
	//endpoint disconnect.{cEndpointId : 1.1.3.1}  {id : 1.1.2.2}, {content : hello world}

}
