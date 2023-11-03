package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	var waitGroup sync.WaitGroup
	input := make(chan interface{})

	for _, command := range cmds {
		output := make(chan interface{})
		waitGroup.Add(1)

		go func(command cmd, input, output chan interface{}) {
			defer waitGroup.Done()
			defer close(output)
			command(input, output)
		}(command, input, output)

		input = output
	}

	waitGroup.Wait()
}

func SelectUsers(in, out chan interface{}) {
	var waitGroup sync.WaitGroup
	processedUsers := make(map[string]bool)
	var mutex sync.Mutex

	for email := range in {
		waitGroup.Add(1)
		go func(email string) {
			defer waitGroup.Done()
			user := GetUser(email)

			if usersAliases[user.Email] != "" {
				user = GetUser(usersAliases[user.Email])
			}

			mutex.Lock()
			if processedUsers[user.Email] {
				mutex.Unlock()
				return
			}

			processedUsers[user.Email] = true
			mutex.Unlock()

			out <- user
		}(email.(string))
	}

	waitGroup.Wait()
}

func processGetMessages(users []User, output chan interface{}, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	messages, _ := GetMessages(users...)
	for _, message := range messages {
		output <- message
	}
}

func SelectMessages(in, out chan interface{}) {
	var waitGroup sync.WaitGroup
	usersBatch := make([]User, 0, GetMessagesMaxUsersBatch)

	for user := range in {
		usersBatch = append(usersBatch, user.(User))

		if len(usersBatch) == GetMessagesMaxUsersBatch {
			waitGroup.Add(1)
			go processGetMessages(usersBatch, out, &waitGroup)
			usersBatch = make([]User, 0, GetMessagesMaxUsersBatch)
		}
	}

	if len(usersBatch) > 0 {
		waitGroup.Add(1)
		go processGetMessages(usersBatch, out, &waitGroup)
	}

	waitGroup.Wait()
}

func CheckSpam(in, out chan interface{}) {
	var waitGroup sync.WaitGroup
	semaphore := make(chan struct{}, HasSpamMaxAsyncRequests)

	for message := range in {
		waitGroup.Add(1)
		semaphore <- struct{}{}

		go func(message MsgID) {
			defer waitGroup.Done()
			hasSpam, _ := HasSpam(message)
			out <- MsgData{
				ID:      message,
				HasSpam: hasSpam,
			}

			<-semaphore
		}(message.(MsgID))
	}

	waitGroup.Wait()
}

func CombineResults(in, out chan interface{}) {
	var waitGroup sync.WaitGroup
	results := make([]MsgData, 0)

	for data := range in {
		results = append(results, data.(MsgData))
	}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		sort.Slice(results, func(i, j int) bool {
			if results[i].HasSpam != results[j].HasSpam {
				return results[i].HasSpam
			}

			return results[i].ID < results[j].ID
		})

		for _, result := range results {
			out <- fmt.Sprintf("%v %d", result.HasSpam, result.ID)
		}
	}()

	waitGroup.Wait()
}
