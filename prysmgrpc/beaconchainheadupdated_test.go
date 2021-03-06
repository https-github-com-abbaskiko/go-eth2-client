// Copyright © 2020 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prysmgrpc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/prysmgrpc"
	"github.com/stretchr/testify/require"
)

func TestBeaconChainheadUpdated(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "Good",
		},
	}

	service, err := prysmgrpc.New(context.Background(),
		prysmgrpc.WithAddress(os.Getenv("PRYSMGRPC_ADDRESS")),
		prysmgrpc.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &testHandler{}
			require.NoError(t, service.AddOnBeaconChainHeadUpdatedHandler(context.Background(), handler))
			ticked := false
			for i := 0; i < 24; i++ {
				if handler.lastSlot != 0 {
					ticked = true
					break
				}
				time.Sleep(time.Second)
			}
			require.True(t, ticked)
		})
	}
}

type testHandler struct {
	lastSlot uint64
}

func (h *testHandler) OnBeaconChainHeadUpdated(ctx context.Context, slot uint64, blockRoot []byte, stateRoot []byte, epochTransition bool) {
	h.lastSlot = slot
}
