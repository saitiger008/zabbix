/*
** Zabbix
** Copyright (C) 2001-2019 Zabbix SIA
**
** This program is free software; you can redistribute it and/or modify
** it under the terms of the GNU General Public License as published by
** the Free Software Foundation; either version 2 of the License, or
** (at your option) any later version.
**
** This program is distributed in the hope that it will be useful,
** but WITHOUT ANY WARRANTY; without even the implied warranty of
** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
** GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License
** along with this program; if not, write to the Free Software
** Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
**/

package systemrun

import (
	"fmt"
	"time"
	"zabbix/internal/agent"
	"zabbix/pkg/plugin"
	"zabbix/pkg/zbxcmd"
)

// Plugin -
type Plugin struct {
	plugin.Base
	enableRemoteCommands int
}

var impl Plugin

func (p *Plugin) Configure(options map[string]string) {
	p.enableRemoteCommands = agent.Options.EnableRemoteCommands
}

// Export -
func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	if p.enableRemoteCommands != 1 {
		return nil, fmt.Errorf("Remote commands are not enabled.")
	}

	if len(params) > 2 {
		return nil, fmt.Errorf("Too many parameters.")
	}

	if len(params) == 0 || len(params[0]) == 0 {
		return nil, fmt.Errorf("Invalid first parameter.")
	}

	if agent.Options.LogRemoteCommands == 1 {
		p.Warningf("Executing command:'%s'", params[0])
	} else {
		p.Debugf("Executing command:'%s'", params[0])
	}

	if len(params) == 1 || params[1] == "" || params[1] == "wait" {
		stdoutStderr, err := zbxcmd.Execute(params[0], time.Second*time.Duration(agent.Options.Timeout))
		if err != nil {
			return nil, err
		}

		p.Debugf("command:'%s' length:%d output:'%.20s'", params[0], len(stdoutStderr), stdoutStderr)

		return stdoutStderr, nil
	} else if params[1] == "nowait" {
		err := zbxcmd.ExecuteBackground(params[0])

		if err != nil {
			return nil, err
		}

		return 1, nil
	}

	return nil, fmt.Errorf("Invalid second parameter.")
}

func init() {
	plugin.RegisterMetrics(&impl, "SystemRun", "system.run", "Run specified command.")
}
