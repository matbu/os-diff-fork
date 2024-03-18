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
package servicecfg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/openstack-k8s-operators/os-diff/pkg/common"
	"github.com/openstack-k8s-operators/os-diff/pkg/godiff"
)

type OpenStackDataPlaneNodeSet struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		NetworkAttachments []string `yaml:"networkAttachments"`
		PreProvisioned     bool     `yaml:"preProvisioned"`
		Services           []string `yaml:"services"`
		Env                []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"env"`
		Nodes struct {
			Standalone struct {
				HostName string `yaml:"hostName"`
				Ansible  struct {
					AnsibleHost string `yaml:"ansibleHost"`
				} `yaml:"ansible"`
				Networks []struct {
					DefaultRoute bool   `yaml:"defaultRoute,omitempty"`
					FixedIP      string `yaml:"fixedIP,omitempty"`
					Name         string `yaml:"name"`
					SubnetName   string `yaml:"subnetName"`
				} `yaml:"networks"`
			} `yaml:"standalone"`
		} `yaml:"nodes"`
		NodeTemplate struct {
			AnsibleSSHPrivateKeySecret string `yaml:"ansibleSSHPrivateKeySecret"`
			ManagementNetwork          string `yaml:"managementNetwork"`
			Ansible                    struct {
				AnsibleUser string `yaml:"ansibleUser"`
				AnsiblePort int    `yaml:"ansiblePort"`
				AnsibleVars struct {
					ServiceNetMap struct {
						NovaAPINetwork     string `yaml:"nova_api_network"`
						NovaLibvirtNetwork string `yaml:"nova_libvirt_network"`
					} `yaml:"service_net_map"`
					EdpmNetworkConfigOverride          string   `yaml:"edpm_network_config_override"`
					EdpmNetworkConfigTemplate          string   `yaml:"edpm_network_config_template"`
					EdpmNetworkConfigHideSensitiveLogs bool     `yaml:"edpm_network_config_hide_sensitive_logs"`
					NeutronPhysicalBridgeName          string   `yaml:"neutron_physical_bridge_name"`
					NeutronPublicInterfaceName         string   `yaml:"neutron_public_interface_name"`
					RoleNetworks                       []string `yaml:"role_networks"`
					NetworksLower                      struct {
						External    string `yaml:"External"`
						InternalAPI string `yaml:"InternalApi"`
						Storage     string `yaml:"Storage"`
						Tenant      string `yaml:"Tenant"`
					} `yaml:"networks_lower"`
					EdpmNodesValidationValidateControllersIcmp bool     `yaml:"edpm_nodes_validation_validate_controllers_icmp"`
					EdpmNodesValidationValidateGatewayIcmp     bool     `yaml:"edpm_nodes_validation_validate_gateway_icmp"`
					EdpmOvnBridgeMappings                      []string `yaml:"edpm_ovn_bridge_mappings"`
					EdpmOvnBridge                              string   `yaml:"edpm_ovn_bridge"`
					EdpmOvnEncapType                           string   `yaml:"edpm_ovn_encap_type"`
					OvnMatchNorthdVersion                      bool     `yaml:"ovn_match_northd_version"`
					OvnMonitorAll                              bool     `yaml:"ovn_monitor_all"`
					EdpmOvnRemoteProbeInterval                 int      `yaml:"edpm_ovn_remote_probe_interval"`
					EdpmOvnOfctrlWaitBeforeClear               int      `yaml:"edpm_ovn_ofctrl_wait_before_clear"`
					TimesyncNtpServers                         []struct {
						Hostname string `yaml:"hostname"`
					} `yaml:"timesync_ntp_servers"`
					EdpmOvnControllerAgentImage   string   `yaml:"edpm_ovn_controller_agent_image"`
					EdpmIscsidImage               string   `yaml:"edpm_iscsid_image"`
					EdpmLogrotateCrondImage       string   `yaml:"edpm_logrotate_crond_image"`
					EdpmNovaComputeContainerImage string   `yaml:"edpm_nova_compute_container_image"`
					EdpmNovaLibvirtContainerImage string   `yaml:"edpm_nova_libvirt_container_image"`
					EdpmOvnMetadataAgentImage     string   `yaml:"edpm_ovn_metadata_agent_image"`
					GatherFacts                   bool     `yaml:"gather_facts"`
					EnableDebug                   bool     `yaml:"enable_debug"`
					EdpmSshdConfigureFirewall     bool     `yaml:"edpm_sshd_configure_firewall"`
					EdpmSshdAllowedRanges         []string `yaml:"edpm_sshd_allowed_ranges"`
					EdpmSelinuxMode               string   `yaml:"edpm_selinux_mode"`
					Plan                          string   `yaml:"plan"`
					EdpmOvsPackages               []string `yaml:"edpm_ovs_packages"`
				} `yaml:"ansibleVars"`
			} `yaml:"ansible"`
		} `yaml:"nodeTemplate"`
	} `yaml:"spec"`
}

func LoadOvsExternalIds(ovsConfig string) map[string]string {
	src, err := ioutil.ReadFile(ovsConfig)
	if err != nil {
		fmt.Println(err)
	}
	data := string(src)
	srcMap := make(map[string]string)

	keyValues := strings.Split(strings.Trim(data, "{}"), ",")
	for _, kv := range keyValues {
		parts := strings.Split(kv, "=")
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(parts[1], "\"")
		srcMap[key] = value
	}
	return srcMap
}

func CompareMappingConfig(srcMap map[string]string, configMapping map[string]string, edpmStruct OpenStackDataPlaneNodeSet) error {
	var report []string
	var msg string
	for k, v := range configMapping {
		value := common.GetNestedFieldValue(edpmStruct.Spec.NodeTemplate.Ansible.AnsibleVars, common.SnakeToCamel(v))
		if srcMap[k] != common.ConvertToString(value) {
			msg = fmt.Sprintf("-%s=%s\n", k, srcMap[k])
			report = append(report, msg)
			msg = fmt.Sprintf("+%s=%s\n", v, common.ConvertToString(value))
			report = append(report, msg)
		}
	}
	godiff.PrintReport(report)
	return nil

}
