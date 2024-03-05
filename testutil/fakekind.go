package testutil

import (
	"context"
	"io"
	"io/ioutil"
	"strings"

	"sigs.k8s.io/kind/pkg/cluster/nodes"
	"sigs.k8s.io/kind/pkg/exec"
)

// FakeProvider
type FakeProvider struct {
	Nodes []nodes.Node
}

func (f *FakeProvider) ListInternalNodes(name string) ([]nodes.Node, error) {
	return f.Nodes, nil
}

type FakeNode struct {
	Cmds []*fakeCmd
	Err  error
}

func (f *FakeNode) CommandContext(ctx context.Context, cmd string, args ...string) exec.Cmd {
	command := &fakeCmd{
		Cmd: strings.Join(append([]string{cmd}, args...), " "),
		err: f.Err,
	}
	f.Cmds = append(f.Cmds, command)
	return command
}

func (f *FakeNode) String() string {
	return "test"
}

// The following functions are not used by our code at all.
func (f *FakeNode) Command(string, ...string) exec.Cmd        { return nil }
func (f *FakeNode) Role() (string, error)                     { return "", nil }
func (f *FakeNode) IP() (ipv4 string, ipv6 string, err error) { return "", "", nil }
func (f *FakeNode) SerialLogs(writer io.Writer) error         { return nil }

type fakeCmd struct {
	Cmd   string
	err   error
	stdin io.Reader
}

func (f *fakeCmd) Run() error {
	if f.stdin != nil {
		// Consume the entire stdin to move the image publish forward.
		ioutil.ReadAll(f.stdin)
	}
	return f.err
}

func (f *fakeCmd) SetStdin(stdin io.Reader) exec.Cmd {
	f.stdin = stdin
	return f
}

// The following functions are not used by our code at all.
func (f *fakeCmd) SetEnv(...string) exec.Cmd    { return f }
func (f *fakeCmd) SetStdout(io.Writer) exec.Cmd { return f }
func (f *fakeCmd) SetStderr(io.Writer) exec.Cmd { return f }
