package main

import (
	"fmt"
	"k8s.io/client-go/util/workqueue"
	"sync"
	"time"
)

func main() {
	// Create queues for virtual server and namespace controllers
	vsQueue := workqueue.NewNamedRateLimitingQueue(
		workqueue.DefaultControllerRateLimiter(), "virtual-server-controller")
	nsQueue := workqueue.NewNamedRateLimitingQueue(
		workqueue.DefaultControllerRateLimiter(), "namespace-controller")

	// Use a WaitGroup to wait for all worker goroutines to finish
	var wg sync.WaitGroup

	// Simulate adding items to the queues by a controller
	wg.Add(1)
	go controllerLoop("virtual-server", vsQueue, &wg)

	wg.Add(1)
	go controllerLoop("namespace", nsQueue, &wg)

	// Simulate processing items in worker goroutines
	wg.Add(1)
	go worker(vsQueue, &wg, "virtual-server")
	wg.Add(1)
	go worker(nsQueue, &wg, "namespace")

	// Allow some time for processing
	time.Sleep(time.Second * 5)

	// Signal controllers to stop and wait for workers to finish
	vsQueue.ShutDown()
	nsQueue.ShutDown()
	wg.Wait()
}

// Simulate a controller loop that adds items to the queue
func controllerLoop(controllerName string, queue workqueue.RateLimitingInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		item := fmt.Sprintf("%s-item-%d", controllerName, i)
		fmt.Printf("Controller %s: Enqueuing item %s\n", controllerName, item)
		queue.Add(item)
		time.Sleep(time.Second)
	}
}

// Simulate a worker goroutine that processes items from the queue
func worker(queue workqueue.RateLimitingInterface, wg *sync.WaitGroup, workerName string) {
	defer wg.Done()
	for {
		item, quit := queue.Get()

		if quit {
			fmt.Printf("Worker %s: Shutting down.\n", workerName)
			return
		}

		// Simulate processing the item
		fmt.Printf("Worker %s: Processing item %s\n", workerName, item)

		// Mark the item as done
		queue.Done(item)
	}
}

