// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var supportedTests = map[string]bool{
	// "difficulty.json":          true, // Testing ETH mainnet config
	"difficultyFrontier.json":     true,
	"difficultyHomestead.json":    true,
	"difficultyByzantium.json":    true,
	"difficultyETC_Atlantis.json": true, // not really filename, but fits pattern
	"difficultyETC_Agharta.json":  true, // "
}

func TestETHDifficultyNDJSON(t *testing.T) {
	filename := filepath.Join(ethBasicTestDir, "mgen_difficulty.ndjson")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		test := DifficultyTest{}
		err = json.Unmarshal(scanner.Bytes(), &test)
		if err != nil {
			t.Fatal(err)
		}

		cfgName := filepath.Dir(test.Name)+".json" // see 'not really filename' comment above
		cfg, ok := ChainConfigs[cfgName]
		if !ok {
			t.Log("Skipping, no config", test.Name, cfgName)
			continue
		}
		t.Run(test.Name, func(t *testing.T) {
			t.Log("Running", test.Name)
			if err := test.runDifficulty(t, &cfg); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestETHDifficulty(t *testing.T) {
	fileNames, _ := filepath.Glob(filepath.Join(ethBasicTestDir, "*"))

	// Loop through each file
	for _, fn := range fileNames {
		fileName := filepath.Base(fn)

		if !supportedTests[fileName] {
			continue
		}

		t.Run(fileName, func(t *testing.T) {
			config := ChainConfigs[fileName]
			tests := make(map[string]DifficultyTest)

			if err := readJsonFile(fn, &tests); err != nil {
				t.Error(err)
			}

			// Loop through each test in file
			for key, test := range tests {
				// Subtest within the JSON file
				t.Run(key, func(t *testing.T) {
					if err := test.runDifficulty(t, &config); err != nil {
						t.Error(err)
					}
				})

			}
		})
	}
}
