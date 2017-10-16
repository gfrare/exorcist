package godEye

import (
	"log"

	"github.com/gfrare/exorcist/rituals"

	"github.com/fsnotify/fsnotify"
	"github.com/gfrare/exorcist/salms"
)

// Watch public function
func Watch(page string) {
	salms.InitAndExecuteMetrics(page)

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
					log.Println("Grimoire has been modified")

					salms.InitAndExecuteMetrics(page)
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(rituals.ConfigurationFile)
	if err != nil {
		log.Fatal(err)
	}
	<-done

}
