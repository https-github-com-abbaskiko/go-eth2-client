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

package tekuhttp_test

import (
	"context"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/tekuhttp"
	"github.com/stretchr/testify/require"
)

type testValidatorIDProvider struct {
	index  uint64
	pubKey string
}

func (t *testValidatorIDProvider) Index(ctx context.Context) (uint64, error) {
	return t.index, nil
}

func (t *testValidatorIDProvider) PubKey(ctx context.Context) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(t.pubKey, "0x"))
}

func TestValidatorBalances(t *testing.T) {
	tests := []struct {
		name       string
		stateID    string
		validators []client.ValidatorIDProvider
	}{
		{
			name:    "Single",
			stateID: "head",
			validators: []client.ValidatorIDProvider{
				&testValidatorIDProvider{
					index:  1000,
					pubKey: "0xb2007d1354db791b924fd35a6b0a8525266a021765b54641f4d415daa50c511204d6acc213a23468f2173e60cc950e26",
				},
			},
		},
		{
			name:    "All",
			stateID: "head",
		},
	}

	service, err := tekuhttp.New(context.Background(),
		tekuhttp.WithAddress(os.Getenv("TEKUHTTP_ADDRESS")),
		tekuhttp.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			balances, err := service.ValidatorBalances(context.Background(), test.stateID, nil)
			require.NoError(t, err)
			require.NotNil(t, balances)
		})
	}
}
