package godEye

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/gfrare/exorcist/salms"
)

// Watch public function
func Watch() {
	salms.InitAndExecuteMetrics()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)

					salms.InitAndExecuteMetrics()
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add("grimoire.json")
	if err != nil {
		log.Fatal(err)
	}
	<-done

}
