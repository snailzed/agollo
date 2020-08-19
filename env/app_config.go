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

package env

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/zouyx/agollo/v4/env/config"
	jsonConfig "github.com/zouyx/agollo/v4/env/config/json"
	"github.com/zouyx/agollo/v4/utils"
)

const (
	appConfigFile     = "app.properties"
	appConfigFilePath = "AGOLLO_CONF"

	defaultNotificationID = int64(-1)
	comma                 = ","

	defaultCluster   = "default"
	defaultNamespace = "application"
)

var (
	//appconfig
	appConfig *config.AppConfig
)

//InitFileConfig 使用文件初始化配置
func InitFileConfig() *config.AppConfig {
	// default use application.properties
	if initConfig, err := InitConfig(nil); err != nil {
		return initConfig
	}
	return nil
}

//InitConfig 使用指定配置初始化配置
func InitConfig(loadAppConfig func() (*config.AppConfig, error)) (*config.AppConfig, error) {
	//init config file
	return getLoadAppConfig(loadAppConfig)
}

//SplitNamespaces 根据namespace字符串分割后，并执行callback函数
func SplitNamespaces(namespacesStr string, callback func(namespace string)) sync.Map {
	namespaces := sync.Map{}
	split := strings.Split(namespacesStr, comma)
	for _, namespace := range split {
		if callback != nil {
			callback(namespace)
		}
		namespaces.Store(namespace, defaultNotificationID)
	}
	return namespaces
}

// set load app config's function
func getLoadAppConfig(loadAppConfig func() (*config.AppConfig, error)) (*config.AppConfig, error) {
	if loadAppConfig != nil {
		return loadAppConfig()
	}
	configPath := os.Getenv(appConfigFilePath)
	if configPath == "" {
		configPath = appConfigFile
	}
	c, e := GetConfigFileExecutor().Load(configPath, Unmarshal)
	if c == nil {
		return nil, e
	}

	return c.(*config.AppConfig), e
}

//GetAppConfig 获取app配置
func GetAppConfig(newAppConfig *config.AppConfig) *config.AppConfig {
	if newAppConfig != nil {
		return newAppConfig
	}
	return appConfig
}

//GetServicesConfigURL 获取服务器列表url
func GetServicesConfigURL(config *config.AppConfig) string {
	return fmt.Sprintf("%sservices/config?appId=%s&ip=%s",
		config.GetHost(),
		url.QueryEscape(config.AppID),
		utils.GetInternal())
}

//GetPlainAppConfig 获取原始配置
func GetPlainAppConfig() *config.AppConfig {
	return appConfig
}

var executeConfigFileOnce sync.Once
var configFileExecutor config.File

//GetConfigFileExecutor 获取文件执行器
func GetConfigFileExecutor() config.File {
	executeConfigFileOnce.Do(func() {
		configFileExecutor = &jsonConfig.ConfigFile{}
	})
	return configFileExecutor
}

//Unmarshal 反序列化
func Unmarshal(b []byte) (interface{}, error) {
	appConfig := &config.AppConfig{
		Cluster:        defaultCluster,
		NamespaceName:  defaultNamespace,
		IsBackupConfig: true,
	}
	err := json.Unmarshal(b, appConfig)
	if utils.IsNotNil(err) {
		return nil, err
	}

	return appConfig, nil
}
