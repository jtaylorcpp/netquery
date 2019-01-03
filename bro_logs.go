package netquery

import (
	"errors"
	"fmt"
	"log"

	"github.com/hpcloud/tail"
	"github.com/jtaylorcpp/broparser"
)

var supportedLogs = map[string]bool{
	"conn.log": true,
	"dns.log":  true,
}

func getBroTailer(path, logname string) (*tail.Tail, error) {
	if _, ok := supportedLogs[logname]; !ok {
		return nil, errors.New(fmt.Sprintf("BRO LOG NOT SUPPORTED: %s, dir: %s", logname, path))
	}

	switch logname {
	case "conn.log":
		return broparser.FollowConn(fmt.Sprintf("%s/%s", path, logname)), nil
	case "dns.log":
		return broparser.FollowDNS(fmt.Sprintf("%s/%s", path, logname)), nil
	default:
		return nil, errors.New("TAILER FAILED")
	}
}

func GetBroTailers(path string, files []string) map[string]*tail.Tail {
	tailMap := map[string]*tail.Tail{}

	for _, file := range files {
		tailer, err := getBroTailer(path, file)
		if err != nil {
			log.Printf("Error For file %s: %s", file, err.Error())
		} else {
			log.Printf("Watching file: %s/%s\n", path, file)
			tailMap[file] = tailer
		}
	}

	return tailMap
}
