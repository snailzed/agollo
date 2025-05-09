/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extension

import (
	"testing"

	"github.com/snailzed/agollo/v4/env/config"
	. "github.com/tevid/gohamcrest"
)

type TestLoadBalance struct {
}

//Load 负载均衡
func (r *TestLoadBalance) Load(servers map[string]*config.ServerInfo) *config.ServerInfo {
	return nil
}

func TestSetLoadBalance(t *testing.T) {
	SetLoadBalance(&TestLoadBalance{})

	balance := GetLoadBalance()

	b := balance.(*TestLoadBalance)
	Assert(t, b, NotNilVal())
}
