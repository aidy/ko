// Copyright 2020 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kind

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/ko/testutil"
	"sigs.k8s.io/kind/pkg/cluster/nodes"
)

func TestWrite(t *testing.T) {
	ctx := context.Background()
	img, err := random.Image(1024, 1)
	if err != nil {
		t.Fatalf("random.Image() = %v", err)
	}

	tag, err := name.NewTag("kind.local/test:new")
	if err != nil {
		t.Fatalf("name.NewTag() = %v", err)
	}

	n1 := &testutil.FakeNode{}
	n2 := &testutil.FakeNode{}
	GetProvider = func() Provider {
		return &testutil.FakeProvider{Nodes: []nodes.Node{n1, n2}}
	}

	if err := Write(ctx, tag, img); err != nil {
		t.Fatalf("Write() = %v", err)
	}

	// Verify the respective command is executed on each node.
	for _, n := range []*testutil.FakeNode{n1, n2} {
		if got, want := len(n.Cmds), 1; got != want {
			t.Fatalf("len(n.cmds) = %d, want %d", got, want)
		}
		c := n.Cmds[0]

		if got, want := c.Cmd, "ctr --namespace=k8s.io images import -"; got != want {
			t.Fatalf("c.cmd = %s, want %s", got, want)
		}
	}
}

func TestTag(t *testing.T) {
	ctx := context.Background()
	oldTag, err := name.NewTag("kind.local/test:test")
	if err != nil {
		t.Fatalf("name.NewTag() = %v", err)
	}

	newTag, err := name.NewTag("kind.local/test:new")
	if err != nil {
		t.Fatalf("name.NewTag() = %v", err)
	}

	n1 := &testutil.FakeNode{}
	n2 := &testutil.FakeNode{}
	GetProvider = func() Provider {
		return &testutil.FakeProvider{Nodes: []nodes.Node{n1, n2}}
	}

	if err := Tag(ctx, oldTag, newTag); err != nil {
		t.Fatalf("Tag() = %v", err)
	}

	// Verify the respective command is executed on each node.
	for _, n := range []*testutil.FakeNode{n1, n2} {
		if got, want := len(n.Cmds), 1; got != want {
			t.Fatalf("len(n.cmds) = %d, want %d", got, want)
		}
		c := n.Cmds[0]

		if got, want := c.Cmd, fmt.Sprintf("ctr --namespace=k8s.io images tag --force %s %s", oldTag, newTag); got != want {
			t.Fatalf("c.cmd = %s, want %s", got, want)
		}
	}
}

func TestFailWithNoNodes(t *testing.T) {
	ctx := context.Background()
	img, err := random.Image(1024, 1)
	if err != nil {
		panic(err)
	}

	oldTag, err := name.NewTag("kind.local/test:test")
	if err != nil {
		t.Fatalf("name.NewTag() = %v", err)
	}

	newTag, err := name.NewTag("kind.local/test:new")
	if err != nil {
		t.Fatalf("name.NewTag() = %v", err)
	}

	GetProvider = func() Provider {
		return &testutil.FakeProvider{}
	}

	if err := Write(ctx, newTag, img); err == nil {
		t.Fatal("Write() = nil, wanted an error")
	}
	if err := Tag(ctx, oldTag, newTag); err == nil {
		t.Fatal("Tag() = nil, wanted an error")
	}
}

func TestFailCommands(t *testing.T) {
	ctx := context.Background()
	img, err := random.Image(1024, 1)
	if err != nil {
		panic(err)
	}

	oldTag, err := name.NewTag("kind.local/test:test")
	if err != nil {
		t.Fatalf("name.NewTag() = %v", err)
	}

	newTag, err := name.NewTag("kind.local/test:new")
	if err != nil {
		t.Fatalf("name.NewTag() = %v", err)
	}

	errTest := errors.New("test")

	n1 := &testutil.FakeNode{Err: errTest}
	n2 := &testutil.FakeNode{Err: errTest}
	GetProvider = func() Provider {
		return &testutil.FakeProvider{Nodes: []nodes.Node{n1, n2}}
	}

	if err := Write(ctx, newTag, img); !errors.Is(err, errTest) {
		t.Fatalf("Write() = %v, want %v", err, errTest)
	}
	if err := Tag(ctx, oldTag, newTag); !errors.Is(err, errTest) {
		t.Fatalf("Write() = %v, want %v", err, errTest)
	}
}
