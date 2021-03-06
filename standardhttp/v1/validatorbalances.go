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

package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type validatorBalancesJSON struct {
	Data []*api.ValidatorBalance `json:"data"`
}

// ValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators are supplied no filter
// will be applied.
func (s *Service) ValidatorBalances(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]spec.Gwei, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/validator_balances", stateID)
	if len(validatorIndices) != 0 {
		ids := make([]string, len(validatorIndices))
		for i := range validatorIndices {
			ids[i] = fmt.Sprintf("%d", validatorIndices[i])
		}
		url = fmt.Sprintf("%s?id=%s", url, strings.Join(ids, "&id="))
	}

	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request validator balances")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain validator balances")
	}

	var validatorBalancesJSON validatorBalancesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&validatorBalancesJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse validator balances")
	}
	if validatorBalancesJSON.Data == nil {
		return nil, errors.New("no validator balances returned")
	}

	res := make(map[spec.ValidatorIndex]spec.Gwei)
	for _, validatorBalance := range validatorBalancesJSON.Data {
		res[validatorBalance.Index] = validatorBalance.Balance
	}
	return res, nil
}
