// SPDX-License-Identifier: AGPL-3.0-or-later
// discord-mpd-bot
// Copyright (C) 2025 skoved
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	client, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not connect to MPD. Check to make sure that MPD is running:", err.Error())
		os.Exit(1)
	}
	statusAttrs, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not get the state of MPD:", err.Error())
		os.Exit(1)
	}
	fmt.Println("State:", statusAttrs["state"])
	if len(os.Args) >= 1 {
		dl := ytdlp.New()
		dl.GetURL()
		dl.AudioFormat("opus")
		dl.AudioQuality("0")
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		res, err := dl.Run(ctx, os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not get song:", os.Args[1])
			os.Exit(1)
		}
		song := strings.Split(res.Stdout, "\n")
		fmt.Println("next song:", song[len(song)-1])
		err = client.Add(song[len(song)-1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not queue song", song[len(song)-1]+":", err.Error())
			os.Exit(1)
		}
	}
	switch statusAttrs["state"] {
	case "pause":
		err = client.Pause(false)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not unpause MPD:", err.Error())
			os.Exit(1)
		}
	case "stop":
		err = client.Play(-1)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not play MPD:", err.Error())
			os.Exit(1)
		}
	}
}
