/*
Copyright 2018 Linode

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"

	"github.com/displague/csi-linode/driver"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "unix:///var/lib/kubelet/plugins/com.linode.csi.linodebs/csi.sock", "CSI endpoint")
		token    = flag.String("token", "", "Linode access token")
		url      = flag.String("url", "https://api.linode.com/", "Linode API URL")
	)
	flag.Parse()

	drv, err := driver.NewDriver(*endpoint, *token, *url)
	if err != nil {
		log.Fatalln(err)
	}

	if err := drv.Run(); err != nil {
		log.Fatalln(err)
	}
}
