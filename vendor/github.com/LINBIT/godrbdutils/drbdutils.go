package godrbdutils

//go:generate stringer -type=Cmd
//go:generate stringer -type=Action

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Cmd int
type Action int

const (
	Drbdsetup Cmd = iota
	Drbdadm
)

const (
	Up Action = iota
	Down
	Adjust
	Attach
	Detach
	Connect
	Disconnect
	Primary
	Secondary
	Create_md
)

type DrbdCmd struct {
	cmd     Cmd
	action  Action
	res     []string
	arg     []string
	timeout time.Duration
}

type DrbdAdm struct {
	res     []string
	timeout time.Duration
}

func NewDrbdAdm(res []string) *DrbdAdm {
	return &DrbdAdm{res: res}
}

func (a *DrbdAdm) Up(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Up, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Down(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Down, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Adjust(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Adjust, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Attach(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Attach, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Detach(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Detach, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Connect(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Connect, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Disconnect(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Disconnect, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Primary(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Primary, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) Secondary(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Secondary, a.res, a.timeout, arg...)
}
func (a *DrbdAdm) CreateMetaData(arg ...string) ([]byte, error) {
	return utilExec(Drbdadm, Create_md, a.res, a.timeout, arg...)
}

func (a *DrbdAdm) SetTimeout(timeout time.Duration) {
	a.timeout = timeout
}

func NewDrbdCmd(cmd Cmd, action Action, res []string, arg ...string) (*DrbdCmd, error) {
	c := DrbdCmd{
		cmd:    cmd,
		action: action,
		res:    res,
		arg:    []string{},
	}
	c.arg = append(c.arg, arg...)
	return &c, nil
}

func (c *DrbdCmd) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *DrbdCmd) CombinedOutput() ([]byte, error) {
	if c.timeout == 0 {
		return c.combinedOutput(nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.combinedOutput(ctx)
}

func (c *DrbdCmd) String() string {
	return strings.Join(c.cmdSlice(), " ")
}

func (c *DrbdCmd) cmdSlice() []string {
	var s = []string{
		strings.ToLower(c.cmd.String()),
		strings.Replace(strings.ToLower(c.action.String()), "_", "-", -1),
	}
	s = append(s, c.arg...)
	for _, r := range c.res {
		s = append(s, r)
	}
	return s
}

func (c *DrbdCmd) combinedOutput(ctx context.Context) ([]byte, error) {
	argv := c.cmdSlice()
	if len(argv) < 2 {
		return nil, fmt.Errorf("Command %v too short", argv)
	}

	var cmd *exec.Cmd
	if ctx != nil {
		cmd = exec.CommandContext(ctx, argv[0], argv[1:]...)
	} else {
		cmd = exec.Command(argv[0], argv[1:]...)
	}

	return cmd.CombinedOutput()
}

func utilExec(cmd Cmd, action Action, res []string, to time.Duration, arg ...string) ([]byte, error) {
	c, err := NewDrbdCmd(cmd, action, res, arg...)
	if err != nil {
		return nil, err
	}
	c.SetTimeout(to)
	return c.CombinedOutput()
}
