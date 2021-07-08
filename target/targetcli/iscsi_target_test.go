package targetcli

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParseCreateIscsiTargetResult(t *testing.T) {
	raw := `Created target iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.64cc17ed0de5.
Created TPG 1.
Global pref auto_add_default_portal=true
Created default portal listening on all IPs (0.0.0.0), port 3260.`

	targetName, tpgId, err := ParseCreateIscsiTargetCmdResult(raw)
	assert.Nil(t, err)
	assert.NotEmpty(t, targetName)
	assert.GreaterOrEqual(t, tpgId, 0)
	spew.Dump(targetName, tpgId, err)
}

func TestParseAddIscsiLunResult(t *testing.T) {
	raw := `Created LUN 0.`

	lunId, err := ParseAddIscsiLunCmdResult(raw)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, lunId, 0)
	spew.Dump(lunId, err)
}
