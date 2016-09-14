package rays

import (
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// New generates new unique RayParam
// Format:
// XXXXXXX-YYYY-ZZZZ-AAAA-BBBB
//
// Where:
//  - XXXXXXX - hostname hash
//  - YYYY - PID
//  - ZZZZ - Startup time
//  - AAAA - Current time
//  - BBBB - Some random value
func New() RayID {
	// Incrementing rays counter with mutex to preserve atomicity
	mutex.Lock()
	incChunk++
	loc := incChunk
	mutex.Unlock()

	// Current time with milliseconds precision
	now := time.Now().UnixNano() / 1000000

	return RayID(fmt.Sprintf("%s-%x-%04x", prefix, now, loc%16000))
}

// prefix contains runtime-static prefix for rays
var prefix string

// incChunk is incremented each time RayId is generated
var incChunk int

// mutex used for incChunk incrementing
var mutex sync.Mutex

// PID contains PID param
var PID PIDParam

// Host contains param for host
var Host HostParam

// InstanceID contains RayID, that should be unchanged
// during program execution
var InstanceID InstanceIDParam

// SupervisorIndex contains supervisord incremental worker index
var SupervisorIndex = SupervisorIndexParam(-1)

func init() {
	pid := os.Getpid()
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}
	Host = HostParam(host)
	PID = PIDParam(pid)

	h := md5.New()
	h.Write([]byte(host))
	sum := h.Sum(nil)[0:3]

	startedAt := time.Now().Unix()

	prefix = fmt.Sprintf("%x-%x-%x", sum, pid, startedAt)

	InstanceID = InstanceIDParam(New().String())

	svGroup := os.Getenv("SUPERVISOR_GROUP_NAME")
	svName := os.Getenv("SUPERVISOR_PROCESS_NAME")
	if svGroup != "" && svName != "" && len(svName) > len(svGroup) {
		if id, err := strconv.Atoi(svName[len(svGroup)+1:]); err == nil {
			SupervisorIndex = SupervisorIndexParam(id)
		}
	}
}
