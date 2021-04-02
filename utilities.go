package main

import "sync"

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func canConsumePage(page string, visitedPages *[]string, mutex *sync.Mutex) bool {
	mutex.Lock()
	_, found := find(*visitedPages, page)
	if !found {
		*visitedPages = append(*visitedPages, page)
	}
	mutex.Unlock()
	return !found
}