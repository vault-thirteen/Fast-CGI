package sr

import (
	"sync"

	cl "github.com/vault-thirteen/Fast-CGI/pkg/Client"
	nvpair "github.com/vault-thirteen/Fast-CGI/pkg/models/NameValuePair"
	pm "github.com/vault-thirteen/Fast-CGI/pkg/models/php"
)

type ScriptRunner struct {
	lock      *sync.Mutex
	requestId uint16
}

func New() (sr *ScriptRunner) {
	return &ScriptRunner{
		lock: new(sync.Mutex),
	}
}

func (sr *ScriptRunner) RunScript(cgiClient *cl.Client, parameters []*nvpair.NameValuePair, stdin []byte) (phpScriptOutput *pm.Data, phpErr error) {
	sr.lock.Lock()
	defer sr.lock.Unlock()

	sr.incRequestId()

	//log.Println(fmt.Sprintf("Request ID = %v.", sr.requestId)) // DEBUG.

	return pm.ExecPhpScriptAndGetHttpData(cgiClient, sr.requestId, parameters, stdin)
}

func (sr *ScriptRunner) incRequestId() {
	sr.requestId++

	// Zero request ID can not be used by user scripts directly.
	if sr.requestId == 0 {
		sr.requestId = 1
	}
}
