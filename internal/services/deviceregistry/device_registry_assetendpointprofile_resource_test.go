package deviceregistry_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AssetEndpointProfileTestResource struct{}

func TestAssetEndpointProfileResource(t *testing.T) {
	// Run all the acceptance tests for the AssetEndpointProfile resource on the cluster.
	// NOTE: this is a combined test rather than separate split out tests due to
	// AssetEndpointProfile resources must be provisioned to the arc-enabled AIO cluster
	// and avoid creating the cluster multiple times.
	testCases := map[string]map[string]func(t *testing.T, randomInteger int){
		// Run the AIO cluster initialization first.
		"InitializeAioCluster": {
			"initializeCluster": testAccAssetEndpointProfile_initializeCluster,
		},
		// Run the acceptance tests for the AssetEndpointProfile resource
		"Resource": {
			"basic":          testAccAssetEndpointProfile_basic,
			"requiresImport": testAccAssetEndpointProfile_requiresImport,
			"completeCertificate":       testAccAssetEndpointProfile_complete_certificate,
			"completeUsernamePassword":  testAccAssetEndpointProfile_complete_usernamePassword,
			"completeAnonymous":         testAccAssetEndpointProfile_complete_anonymous,
			"update":         testAccAssetEndpointProfile_update,
		},
	}

	// Generate a random integer of size 18 that will stay the same in all test cases.
	// This value will be used to create unique names for the AIO cluster's infra resources
	// (such as the resource group name, VM name, etc) but will be kept constant
	// so that the same cluster is used for all the acceptance tests.
	constantRandomInt := acceptance.RandTimeInt()

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t, constantRandomInt)
				})
			}
		})
	}
}

func testAccAssetEndpointProfile_basic(t *testing.T, randomInteger int) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue(""),
				check.That(data.ResourceName).Key("authentication_method").HasValue(""),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func testAccAssetEndpointProfile_complete_certificate(t *testing.T, randomInteger int) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeCertificate(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Certificate"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue("myCertificateRef"),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func testAccAssetEndpointProfile_complete_usernamePassword(t *testing.T, randomInteger int) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeUsernamePassword(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("UsernamePassword"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue("myUsernameRef"),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue("myPasswordRef"),
			),
		},
		data.ImportStep(),
	})
}

func testAccAssetEndpointProfile_complete_anonymous(t *testing.T, randomInteger int) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeAnonymous(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Anonymous"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func testAccAssetEndpointProfile_requiresImport(t *testing.T, randomInteger int) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, randomInteger)
		}),
	})
}

func testAccAssetEndpointProfile_update(t *testing.T, randomInteger int) {
	data := acceptance.BuildTestData(t, "azurerm_device_registry_asset_endpoint_profile", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{ // first create the resource
			Config: r.basic(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to certificate
			Config: r.completeCertificate(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Certificate"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue("myCertificateRef"),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to username/password
			Config: r.completeUsernamePassword(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("UsernamePassword"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue("myUsernameRef"),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue("myPasswordRef"),
			),
		},
		data.ImportStep(),
		{ // update the authentication method to anonymous
			Config: r.completeAnonymous(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("target_address").HasValue("opc.tcp://foo"),
				check.That(data.ResourceName).Key("endpoint_profile_type").HasValue("OpcUa"),
				check.That(data.ResourceName).Key("discovered_asset_endpoint_profile_ref").HasValue("discoveredAssetEndpointProfile123"),
				check.That(data.ResourceName).Key("additional_configuration").HasValue("{\"foo\": \"bar\"}"),
				check.That(data.ResourceName).Key("authentication_method").HasValue("Anonymous"),
				check.That(data.ResourceName).Key("x509_credentials_certificate_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_username_secret_name").HasValue(""),
				check.That(data.ResourceName).Key("username_password_credentials_password_secret_name").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func testAccAssetEndpointProfile_initializeCluster(t *testing.T, randomInteger int) {
	data := acceptance.BuildTestData(t, "azurerm_linux_virtual_machine", "test")
	r := AssetEndpointProfileTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data, randomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (AssetEndpointProfileTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := assetendpointprofiles.ParseAssetEndpointProfileID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DeviceRegistry.AssetEndpointProfileClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r AssetEndpointProfileTestResource) basic(data acceptance.TestData, randomInteger int) string {
	template := r.constantsTemplate(data, randomInteger)
	
	// AssetEndpointProfile can have different random int for the name. Only AIO cluster infra needs to have the same random int.
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name                                  = "acctest-assetendpointprofile-%[2]d"
	resource_group_name                   = azurerm_resource_group.test.name
	extended_location_name                = "${azurerm_resource_group.test.id}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
	extended_location_type                = "CustomLocation"
	target_address                        = "opc.tcp://foo"
	endpoint_profile_type                 = "OpcUa"
	discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
	location                              = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeCertificate(data acceptance.TestData, randomInteger int) string {
	template := r.constantsTemplate(data, randomInteger)

	// AssetEndpointProfile can have different random int for the name. Only AIO cluster infra needs to have the same random int.
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name                                     = "acctest-assetendpointprofile-%[2]d"
	resource_group_name                      = azurerm_resource_group.test.name
	extended_location_name                   = "${azurerm_resource_group.test.id}/providers/Microsoft.ExtendedLocation/customLocations/{local.custom_location}"
	extended_location_type                   = "CustomLocation"
	target_address                           = "opc.tcp://foo"
	endpoint_profile_type                    = "OpcUa"
	discovered_asset_endpoint_profile_ref    = "discoveredAssetEndpointProfile123"
	additional_configuration                 = "{\"foo\": \"bar\"}"
	authentication_method                    = "Certificate"
	x509_credentials_certificate_secret_name = "myCertificateRef"
	location                                 = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeUsernamePassword(data acceptance.TestData, randomInteger int) string {
	template := r.constantsTemplate(data, randomInteger)

	// AssetEndpointProfile can have different random int for the name. Only AIO cluster infra needs to have the same random int.
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name                                               = "acctest-assetendpointprofile-%[2]d"
	resource_group_name                                = azurerm_resource_group.test.name
	extended_location_name                             = "${azurerm_resource_group.test.id}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
	extended_location_type                             = "CustomLocation"
	target_address                                     = "opc.tcp://foo"
	endpoint_profile_type                              = "OpcUa"
	discovered_asset_endpoint_profile_ref              = "discoveredAssetEndpointProfile123"
	additional_configuration                           = "{\"foo\": \"bar\"}"
	authentication_method                              = "UsernamePassword"
	username_password_credentials_username_secret_name = "myUsernameRef"
	username_password_credentials_password_secret_name = "myPasswordRef"
	location                                           = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) completeAnonymous(data acceptance.TestData, randomInteger int) string {
	template := r.constantsTemplate(data, randomInteger)

	// AssetEndpointProfile can have different random int for the name. Only AIO cluster infra needs to have the same random int.
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "test" {
	name                                  = "acctest-assetendpointprofile-%[2]d"
  resource_group_name                   = azurerm_resource_group.test.name
	extended_location_name                = "${azurerm_resource_group.test.id}/providers/Microsoft.ExtendedLocation/customLocations/${local.custom_location}"
	extended_location_type                = "CustomLocation"
	target_address                        = "opc.tcp://foo"
	endpoint_profile_type                 = "OpcUa"
	discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
	additional_configuration              = "{\"foo\": \"bar\"}"
	authentication_method                 = "Anonymous"
	location                              = "%[3]s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AssetEndpointProfileTestResource) requiresImport(data acceptance.TestData, randomInteger int) string {
	template := r.basic(data, randomInteger)
	return fmt.Sprintf(`
%s

resource "azurerm_device_registry_asset_endpoint_profile" "import" {
	name 					                        = azurerm_device_registry_asset_endpoint_profile.test.name
	resource_group_name                   = azurerm_device_registry_asset_endpoint_profile.test.resource_group_name
	extended_location_name                = azurerm_device_registry_asset_endpoint_profile.test.extended_location_name
	extended_location_type                = azurerm_device_registry_asset_endpoint_profile.test.extended_location_type
	target_address                        = azurerm_device_registry_asset_endpoint_profile.test.target_address
	endpoint_profile_type                 = azurerm_device_registry_asset_endpoint_profile.test.endpoint_profile_type
	discovered_asset_endpoint_profile_ref = "discoveredAssetEndpointProfile123"
	location                              = azurerm_device_registry_asset_endpoint_profile.test.location
}

`, template)
}

/*
Creates the terraform template for constants needed for the AIO cluster infra.
*/
func (AssetEndpointProfileTestResource) constantsTemplate(data acceptance.TestData, randomInteger int) string {
	// Trim the random value (from acceptance.RandTimeInt which is 18 digits) to 10 digits
	// to avoid exceeding the maximum length of the storage account name (24 chars max).
	trimmedRandomInteger := randomInteger % 10000000000
	return fmt.Sprintf(`
locals {
	custom_location           = "adr-acctest-cl%[1]d"
	storage_account           = "acctestsa%[2]d"
	schema_registry           = "acctest-sr-%[1]d"
	schema_registry_namespace = "acctests-rn-%[1]d"
	resource_group_name       = "adr-acctest-rg-%[1]d"
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}
`, randomInteger, trimmedRandomInteger)
}

/*
The terraform template for all the resources needed to create an AIO cluster on a VM
which the acceptance tests' AssetEndpointProfile resources will be provisioned to.
*/
func (r AssetEndpointProfileTestResource) template(data acceptance.TestData, randomInteger int) string {
	clientId := os.Getenv("ARM_CLIENT_ID")
	constantsTemplate := r.constantsTemplate(data, randomInteger)
	credential := r.getCredentials()
	provisionTemplate := r.provisionTemplate(data, credential, randomInteger)

	return fmt.Sprintf(`
%[5]s

resource "azurerm_resource_group" "test" {
  name     = local.resource_group_name
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
  }
}

resource "azurerm_network_security_group" "my_terraform_nsg" {
  name                = "myNetworkSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  lifecycle {
    ignore_changes = [
      security_rule,
    ]
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.my_terraform_nsg.id
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctestVM-%[1]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
	//resource_group_name             = "adr-terraform-acctest-rg"
	//location                        = "%[2]s"
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "%[3]s"
  provision_vm_agent              = false
  allow_extension_operations      = false
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]
  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }
  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

	identity {
		type = "SystemAssigned, UserAssigned"
		identity_ids = [
			%[6]s
		]
	}

	%[4]s

  depends_on = [
    azurerm_network_interface_security_group_association.test
  ]
}
`, randomInteger, data.Locations.Primary, credential, provisionTemplate, constantsTemplate, clientId)
}

/*
Copies the scripts and files needed to create and provision the AIO cluster on the VM.
Then ssh's into the VM and executes the cluster setup scripts.
In case of errors during remote execution of scripts, the logs are written to a file `agent_log` on the VM.
*/
func (r AssetEndpointProfileTestResource) provisionTemplate(data acceptance.TestData, credential string, randomInteger int) string {
	// Get client secrets from env vars because we need them 
	// to remote execute az cli commands on the VM.
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	
	return fmt.Sprintf(`
connection {
 type     = "ssh"
 host     = azurerm_public_ip.test.ip_address
 user     = "adminuser"
 password = "%[1]s"
}

provisioner "file" {
 content = templatefile("testdata/setup_aio_cluster.sh.tftpl", {
   subscription_id     = data.azurerm_client_config.current.subscription_id
   resource_group_name = azurerm_resource_group.test.name
   cluster_name        = "acctest-akcc-%[2]d"
   location            = azurerm_resource_group.test.location
	 custom_location     = local.custom_location
	 storage_account     = local.storage_account
	 schema_registry     = local.schema_registry
	 schema_registry_namespace = local.schema_registry_namespace
   tenant_id           = data.azurerm_client_config.current.tenant_id
	 client_id           = "%[4]s"
	 client_secret       = "%[5]s"
   working_dir         = "%[3]s"
 })
 destination = "%[3]s/setup_aio_cluster.sh"
}

provisioner "file" {
 source      = "testdata/setup_aio_cluster.py"
 destination = "%[3]s/setup_aio_cluster.py"
}

provisioner "remote-exec" {
 inline = [
   "sudo sed -i 's/\r$//' %[3]s/setup_aio_cluster.sh",
   "sudo chmod +x %[3]s/setup_aio_cluster.sh",
   "bash %[3]s/setup_aio_cluster.sh > %[3]s/agent_log",
 ]
}
`, credential, data.RandomInteger, "/home/adminuser", clientId, clientSecret)
}

// Generates a random password for the VM.
func (AssetEndpointProfileTestResource) getCredentials() string {
	return fmt.Sprintf("P@$$w0rd%d!", rand.Intn(10000))
}