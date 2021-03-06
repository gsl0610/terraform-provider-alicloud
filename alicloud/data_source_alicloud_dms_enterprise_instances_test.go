package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudDmsEnterprisesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_dms_enterprise_instances.default"
	name := fmt.Sprintf("tf_testAccDmsEnterpriseInstancesDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		name, dataSourceDmsEnterpriseInstancesConfigDependence)

	searchkeyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"search_key": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"search_key": name + "fake",
		}),
	}
	instancealiasRegexConfConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"net_type":             "CLASSIC",
			"instance_type":        "mysql",
			"env_type":             "test",
			"instance_alias_regex": "tf_testAcc",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"net_type":             "CLASSIC",
			"instance_type":        "mysql",
			"env_type":             "test",
			"instance_alias_regex": "tf_testAcc-fake",
		}),
	}
	var existDmsEnterpriseInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":                   "1",
			"instances.0.data_link_name":    "",
			"instances.0.database_password": CHECKSET,
			"instances.0.database_user":     "tftestnormal",
			"instances.0.dba_id":            CHECKSET,
			"instances.0.dba_nick_name":     CHECKSET,
			"instances.0.ddl_online":        "0",
			"instances.0.ecs_instance_id":   "",
			"instances.0.ecs_region":        "cn-hangzhou",
			"instances.0.env_type":          "test",
			"instances.0.export_timeout":    CHECKSET,
			"instances.0.host":              CHECKSET,
			"instances.0.instance_alias":    CHECKSET,
			"instances.0.instance_id":       CHECKSET,
			"instances.0.instance_source":   "RDS",
			"instances.0.instance_type":     "mysql",
			"instances.0.port":              "3306",
			"instances.0.query_timeout":     CHECKSET,
			"instances.0.safe_rule_id":      CHECKSET,
			"instances.0.sid":               "",
			"instances.0.status":            CHECKSET,
			"instances.0.use_dsql":          "0",
			"instances.0.vpc_id":            "",
		}
	}

	var fakeDmsEnterpriseInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
		}
	}

	var DmsEnterpriseInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDmsEnterpriseInstancesMapFunc,
		fakeMapFunc:  fakeDmsEnterpriseInstancesMapFunc,
	}

	DmsEnterpriseInstancesCheckInfo.dataSourceTestCheck(t, rand, searchkeyConf, instancealiasRegexConfConf)
}

func dataSourceDmsEnterpriseInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "creation" {
		  default = "Rds"
		}
		data "alicloud_account" "current"{
		}
		variable "name" {
		  default = "dbconnectionbasic"
		}
	
		data "alicloud_zones" "default" {
		  available_resource_creation = var.creation
		}
	
		resource "alicloud_vpc" "default" {
		  name       = var.name
		  cidr_block = "172.16.0.0/16"
		}
	
		resource "alicloud_vswitch" "default" {
		  vpc_id            = alicloud_vpc.default.id
		  cidr_block        = "172.16.0.0/24"
		  availability_zone = data.alicloud_zones.default.zones[0].id
		  name              = var.name
		}
		resource "alicloud_security_group" "default" {
		  name   = var.name
		  vpc_id = alicloud_vpc.default.id
		}
		resource "alicloud_db_instance" "instance" {
		  engine           = "MySQL"
		  engine_version   = "5.7"
		  instance_type    = "rds.mysql.t1.small"
		  instance_storage = "10"
		  vswitch_id       = alicloud_vswitch.default.id
		  instance_name    = var.name
		  security_ips     = ["100.104.5.0/24","192.168.0.6"]
		}
	
		resource "alicloud_db_account" "account" {
		  instance_id = alicloud_db_instance.instance.id
		  name        = "tftestnormal"
		  password    = "Test12345"
		  type        = "Normal"
		}
	
		resource "alicloud_dms_enterprise_instance" "default" {
		  dba_uid           =  tonumber(data.alicloud_account.current.id)
		  host              =  alicloud_db_instance.instance.connection_string
		  port              =  "3306"
		  network_type      =	 "VPC"
		  safe_rule         =	"自由操作"
		  tid               =  "13429"
		  instance_type     =	 "mysql"
		  instance_source   =	 "RDS"
		  env_type          =	 "test"
		  database_user     =	 alicloud_db_account.account.name
		  database_password =	 alicloud_db_account.account.password
		  instance_alias    =	 %s
		  query_timeout     =	 "70"
		  export_timeout    =	 "2000"
		  ecs_region        =	 "cn-hangzhou"
		  ddl_online        =	 "0"
		  use_dsql          =	 "0"
		  data_link_name    =	 ""
		}
	`, name)
}
