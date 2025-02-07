provider "azurerm" {
	features {}

	subscription_id = "efb15086-3322-405d-a9d0-c35715a9b722"
	tenant_id = "72f988bf-86f1-41af-91ab-2d7cd011db47"
	# client_id = "17977393-e53c-45dd-a8c1-4d1a06c24298"
    # use_msi = true
	resource_provider_registrations = "none"
}

data "azurerm_subscription" "current" {}

/*
resource "azurerm_device_registry_asset_endpoint_profile" "test2" {
  	name             = "myassetendpointprofilebasic2"
	resource_group_name = "adr-terraform-test-113553226"
	extended_location_name = "/subscriptions/efb15086-3322-405d-a9d0-c35715a9b722/resourceGroups/adr-terraform-test-113553226/providers/Microsoft.ExtendedLocation/customLocations/location-2h2vr"
	extended_location_type = "CustomLocation"
	target_address = "opc.tcp://foo2"
	endpoint_profile_type = "OpcUa"
	//authentication_method = "UsernamePassword"
	username_password_credentials_username_secret_name = "myusernameref"
	username_password_credentials_password_secret_name = "mypasswordref" 
	discovered_asset_endpoint_profile_ref = "mydaepref"
	additional_configuration = "{\"foo\": \"bar\"}"
	location         = "westus2"
}



resource "azurerm_device_registry_asset" "testasset" {
    asset_endpoint_profile_ref     = "myaepref"
    attributes                     = {
        "foo" = "bar"
        "x"   = "y"
    }
    default_datasets_configuration = jsonencode(
        {
            defaultPublishingInterval = 200
            defaultQueueSize          = 10
            defaultSamplingInterval   = 500
        }
    )
    default_events_configuration   = jsonencode(
        {
            defaultPublishingInterval = 200
            defaultQueueSize          = 10
            defaultSamplingInterval   = 500
        }
    )
    default_topic_path             = "/path/defaultTopic"
    default_topic_retain           = "Keep"
    description                    = "this is my asset"
    discovered_asset_refs          = [
        "foo",
        "bar",
        "baz",
    ]
    display_name                   = "my asset"
    documentation_uri              = "https://example.com/about"
    enabled                        = false
    extended_location_name         = "/subscriptions/efb15086-3322-405d-a9d0-c35715a9b722/resourceGroups/adr-terraform-test-113553226/providers/Microsoft.ExtendedLocation/customLocations/location-2h2vr"
    extended_location_type         = "CustomLocation"
    external_asset_id              = "foobar"
    hardware_revision              = "1.0"
    location                       = "westus2"
    manufacturer                   = "contoso"
    manufacturer_uri               = "https://example.com"
    model                          = "model1"
    name                           = "myassetbasic"
    product_code                   = "42"
    resource_group_name            = "adr-terraform-test-113553226"
    serial_number                  = "1234"
    software_revision              = "1.0"
    tags                           = {
        "sensor" = "temperature,humidity"
    }

		datasets {
			name = "dataset1"
			dataset_configuration = jsonencode(
				{
					publishingInterval = 7
					queueSize          = 8
					samplingInterval   = 1000
				}
			)
			topic_path = "/path/dataset1"
			topic_retain = "Keep"
			
			data_points {
				name = "datapoint1"
				data_source = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt1"
				observability_mode = "Log"
				data_point_configuration = jsonencode(
					{
						publishingInterval = 7
						queueSize          = 8
						samplingInterval   = 1000
					}
				)
			}

			data_points {
				name = "datapoint2"
				data_source = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt2"
				observability_mode = "None"
				data_point_configuration = jsonencode(
					{
						publishingInterval = 7
						queueSize          = 8
						samplingInterval   = 1000
					}
				)
			}
		}

    events {
        event_configuration = jsonencode(
            {
                publishingInterval = 7
                queueSize          = 8
                samplingInterval   = 1000
            }
        )
        event_notifier      = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt3"
        name                = "event1"
        observability_mode  = "Log"
        topic_path          = "/path/event1"
        topic_retain        = "Never"
    }
    events {
        event_configuration = jsonencode(
            {
                publishingInterval = 7
                queueSize          = 8
                samplingInterval   = 1000
            }
        )
        event_notifier      = "nsu=http://microsoft.com/Opc/OpcPlc/;s=FastUInt4"
        name                = "event2"
        observability_mode  = "None"
        topic_path          = "/path/event2"
        topic_retain        = "Keep"
    }
}

*/