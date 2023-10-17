/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */
package cmd

import (
	"os-diff/pkg/godiff"

	"github.com/spf13/cobra"
)

// Diff parameters
var debug bool
var remote bool
var file1Cmd string
var file2Cmd string

var diffCmd = &cobra.Command{
	Use:   "diff [file1] [file2]",
	Short: "Print diff for two specific files",
	Long: `Print diff for files provided via the command line: For example:
./os-diff diff tests/podman/keystone.conf tests/ocp/keystone.conf
Example for remote diff:
export CMD1="ssh -F ssh.config standalone podman exec a6e1ca049eee"
export CMD2="oc exec glance-external-api-6cf6c98564-blggc -c glance-api --"
./os-diff diff /etc/glance/glance-api.conf /etc/glance/glance.conf.d/00-config.conf --file1-cmd "$CMD1" --file2-cmd "$CMD2" /etc/glance/glance-api.conf -d /etc/glance/glance.conf.d/00-config.conf --remote`,
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		file2 := args[1]
		if remote {
			godiff.CompareFilesFromRemote(file1, file2, file1Cmd, file2Cmd, debug)
		} else {
			godiff.CompareFiles(file1, file2, true, debug)
		}
	},
}

func init() {
	diffCmd.Flags().StringVarP(&file1Cmd, "file1-cmd", "", "", "Remote command for the file1 configuration file.")
	diffCmd.Flags().StringVarP(&file1Cmd, "file2-cmd", "", "", "Remote command for the file2 configuration file.")
	diffCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug.")
	diffCmd.Flags().BoolVar(&remote, "remote", false, "Run the diff remotely.")
	rootCmd.AddCommand(diffCmd)
}
