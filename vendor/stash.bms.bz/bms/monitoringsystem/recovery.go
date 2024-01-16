package monitoringsystem

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func panicRecover() {
	if r := recover(); r != nil {
		err := fmt.Sprintf("%v", r)
		host, _ := os.Hostname()
		data, _ := json.Marshal(map[string]string{
			"host":        host,
			"app":         "monitoring-lib",
			"file":        "",
			"method":      "panicRecover",
			"type":        "fatal",
			"component":   "Application",
			"code":        "MONITORING.LIBRARY.PANIC",
			"description": err,
			"category":    "UnknownError",
			"doc":         "",
			"ts":          time.Now().String(),
		})
		fmt.Println(string(data))
	}
}
