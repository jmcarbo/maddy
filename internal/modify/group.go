/*
Maddy Mail Server - Composable all-in-one email server.
Copyright © 2019-2020 Max Mazurov <fox.cpp@disroot.org>, Maddy Mail Server contributors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package modify

import (
	"context"

	"github.com/emersion/go-message/textproto"
	"github.com/jmcarbo/maddy/framework/buffer"
	"github.com/jmcarbo/maddy/framework/config"
	modconfig "github.com/jmcarbo/maddy/framework/config/module"
	"github.com/jmcarbo/maddy/framework/module"
)

type (
	// Group wraps multiple modifiers and runs them serially.
	//
	// It is also registered as a module under 'modifiers' name and acts as a
	// module group.
	Group struct {
		instName  string
		Modifiers []module.Modifier
	}

	groupState struct {
		states []module.ModifierState
	}
)

func (g *Group) Init(cfg *config.Map) error {
	for _, node := range cfg.Block.Children {
		mod, err := modconfig.MsgModifier(cfg.Globals, append([]string{node.Name}, node.Args...), node)
		if err != nil {
			return err
		}

		g.Modifiers = append(g.Modifiers, mod)
	}

	return nil
}

func (g *Group) Name() string {
	return "modifiers"
}

func (g *Group) InstanceName() string {
	return g.instName
}

func (g Group) ModStateForMsg(ctx context.Context, msgMeta *module.MsgMetadata) (module.ModifierState, error) {
	gs := groupState{}
	for _, modifier := range g.Modifiers {
		state, err := modifier.ModStateForMsg(ctx, msgMeta)
		if err != nil {
			// Free state objects we initialized already.
			for _, state := range gs.states {
				state.Close()
			}
			return nil, err
		}
		gs.states = append(gs.states, state)
	}
	return gs, nil
}

func (gs groupState) RewriteSender(ctx context.Context, mailFrom string) (string, error) {
	var err error
	for _, state := range gs.states {
		mailFrom, err = state.RewriteSender(ctx, mailFrom)
		if err != nil {
			return "", err
		}
	}
	return mailFrom, nil
}

func (gs groupState) RewriteRcpt(ctx context.Context, rcptTo string) (string, error) {
	var err error
	for _, state := range gs.states {
		rcptTo, err = state.RewriteRcpt(ctx, rcptTo)
		if err != nil {
			return "", err
		}
	}
	return rcptTo, nil
}

func (gs groupState) RewriteBody(ctx context.Context, h *textproto.Header, body buffer.Buffer) error {
	for _, state := range gs.states {
		if err := state.RewriteBody(ctx, h, body); err != nil {
			return err
		}
	}
	return nil
}

func (gs groupState) Close() error {
	// We still try close all state objects to minimize
	// resource leaks when Close fails for one object..

	var lastErr error
	for _, state := range gs.states {
		if err := state.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func init() {
	module.Register("modifiers", func(_, instName string, _, _ []string) (module.Module, error) {
		return &Group{
			instName: instName,
		}, nil
	})
}
