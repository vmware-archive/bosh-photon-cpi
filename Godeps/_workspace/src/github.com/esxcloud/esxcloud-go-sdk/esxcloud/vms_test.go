package esxcloud

import (
	"reflect"
	"testing"
)

func TestCreateGetDeleteVM(t *testing.T) {
	// Create tenant
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10)}
	tenantTask, _ := client.Tenants.Create(tenantSpec)

	// Create resource ticket
	resSpec := &ResourceTicketCreateSpec{
		Name:   randomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 4, Key: "vm.memory"}},
	}
	_, _ = client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, resSpec)

	// Create project
	projSpec := &ProjectCreateSpec{
		ResourceTicketReservation{
			resSpec.Name,
			[]QuotaLineItem{
				QuotaLineItem{"GB", 4, "vm.memory"},
				QuotaLineItem{"COUNT", 100, "vm"},
			},
		},
		randomString(10),
	}
	projTask, _ := client.Tenants.CreateProject(tenantTask.Entity.ID, projSpec)

	// Create flavor
	flavorSpec := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "ephemeral-disk.cost"}},
		"ephemeral-disk",
		randomString(10),
	}
	_, _ = client.Flavors.Create(flavorSpec)

	// Upload image
	imagePath := "../testdata/tty_tiny.ova"
	imageTask, _ := client.Images.CreateFromFile(imagePath, nil)
	mockTask = createMockTask("CREATE_IMAGE", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	imageTask, _ = client.Tasks.Wait(imageTask.ID)

	vmFlavorSpec := &FlavorCreateSpec{
		Name: randomString(10),
		Kind: "vm",
		Cost: []QuotaLineItem{
			QuotaLineItem{"GB", 2, "vm.memory"},
			QuotaLineItem{"COUNT", 4, "vm.cpu"},
		},
	}
	_, _ = client.Flavors.Create(vmFlavorSpec)

	// Create VM
	mockTask = createMockTask("CREATE_VM", "QUEUED")
	server.SetResponseJson(200, mockTask)
	vmSpec := &VmCreateSpec{
		Flavor:        vmFlavorSpec.Name,
		SourceImageID: imageTask.Entity.ID,
		AttachedDisks: []AttachedDisk{
			AttachedDisk{
				CapacityGB: 1,
				Flavor:     flavorSpec.Name,
				Kind:       "ephemeral-disk",
				Name:       randomString(10),
				State:      "STARTED",
				BootDisk:   true,
			},
		},
		Name: randomString(10),
	}

	vmCreateTask, err := client.Projects.CreateVM(projTask.Entity.ID, vmSpec)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vmCreateTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if vmCreateTask.Operation != "CREATE_VM" {
		t.Error("Expected task operation to be CREATE_VM")
	}
	if vmCreateTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for VM creation
	mockTask = createMockTask("CREATE_VM", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	vmCreateTask, err = client.Tasks.Wait(vmCreateTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vmCreateTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Set VM metadata
	metadata := &VmMetadata{Metadata: map[string]interface{}{"key1": "value1"}}
	_, err = client.VMs.SetMetadata(vmCreateTask.Entity.ID, metadata)
	if err != nil {
		t.Error("Not expecting error")
	}

	// Get VM
	mockVm := &VM{Name: vmSpec.Name, Metadata: metadata.Metadata}
	server.SetResponseJson(200, mockVm)
	vm, err := client.VMs.Get(vmCreateTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vm == nil {
		t.Error("Not expecting VM to be nil")
	}
	if vm.Name != vmSpec.Name {
		t.Error("Did not see expected VM from Get VM")
	}
	if !reflect.DeepEqual(metadata.Metadata, vm.Metadata) {
		t.Error("VM metadata did not match expected")
	}

	// Cleanup VM
	mockTask = createMockTask("DELETE_VM", "QUEUED")
	server.SetResponseJson(200, mockTask)
	vmDeleteTask, err := client.VMs.Delete(vmCreateTask.Entity.ID, true)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vmDeleteTask.Operation != "DELETE_VM" {
		t.Error("Expected task operation to be DELETE_VM")
	}
	if vmDeleteTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	mockTask = createMockTask("DELETE_VM", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	vmDeleteTask, err = client.Tasks.Wait(vmDeleteTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vmDeleteTask.Operation != "DELETE_VM" {
		t.Error("Expected task operation to be DELETE_VM")
	}
	if vmDeleteTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Cleanup image
	imageTask, err = client.Images.Delete(imageTask.Entity.ID)
	imageTask, err = client.Tasks.Wait(imageTask.ID)
	if err != nil {
		t.Error("Not expecting error when deleting image")
		t.Log(err)
	}

	// Cleanup project
	_, err = client.Projects.Delete(projTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting project")
		t.Log(err)
	}

	// Cleanup tenant
	_, err = client.Tenants.Delete(tenantTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting tenant")
		t.Log(err)
	}
}

func TestAttachDetachDisk(t *testing.T) {
	// Create tenant
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10)}
	tenantTask, _ := client.Tenants.Create(tenantSpec)

	// Create resource ticket
	resSpec := &ResourceTicketCreateSpec{
		Name:   randomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 4, Key: "vm.memory"}},
	}
	_, _ = client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, resSpec)

	// Create project
	projSpec := &ProjectCreateSpec{
		ResourceTicketReservation{
			resSpec.Name,
			[]QuotaLineItem{
				QuotaLineItem{"GB", 4, "vm.memory"},
				QuotaLineItem{"COUNT", 100, "vm"},
			},
		},
		randomString(10),
	}
	projTask, _ := client.Tenants.CreateProject(tenantTask.Entity.ID, projSpec)

	// Create ephemeral disk flavor
	ephemeralFlavor := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "ephemeral-disk.cost"}},
		"ephemeral-disk",
		randomString(10),
	}
	_, _ = client.Flavors.Create(ephemeralFlavor)

	// Create persistent disk flavor
	persistentFlavor := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "persistent-disk.cost"}},
		"persistent-disk",
		randomString(10),
	}
	_, _ = client.Flavors.Create(persistentFlavor)

	// Create persistent disk
	diskSpec := &DiskCreateSpec{
		Flavor:     persistentFlavor.Name,
		Kind:       "persistent-disk",
		CapacityGB: 1,
		Name:       randomString(10),
	}
	diskTask, _ := client.Projects.CreateDisk(projTask.Entity.ID, diskSpec)
	diskTask, _ = client.Tasks.Wait(diskTask.ID)

	// Upload image
	imagePath := "../testdata/tty_tiny.ova"
	imageTask, _ := client.Images.CreateFromFile(imagePath, nil)
	mockTask = createMockTask("CREATE_IMAGE", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	imageTask, _ = client.Tasks.Wait(imageTask.ID)

	vmFlavorSpec := &FlavorCreateSpec{
		Name: randomString(10),
		Kind: "vm",
		Cost: []QuotaLineItem{
			QuotaLineItem{"GB", 2, "vm.memory"},
			QuotaLineItem{"COUNT", 4, "vm.cpu"},
		},
	}
	_, _ = client.Flavors.Create(vmFlavorSpec)

	// Create VM
	mockTask = createMockTask("CREATE_VM", "QUEUED")
	server.SetResponseJson(200, mockTask)
	vmSpec := &VmCreateSpec{
		Flavor:        vmFlavorSpec.Name,
		SourceImageID: imageTask.Entity.ID,
		AttachedDisks: []AttachedDisk{
			AttachedDisk{
				CapacityGB: 1,
				Flavor:     ephemeralFlavor.Name,
				Kind:       "ephemeral-disk",
				Name:       randomString(10),
				State:      "STARTED",
				BootDisk:   true,
			},
		},
		Name: randomString(10),
	}
	vmCreateTask, _ := client.Projects.CreateVM(projTask.Entity.ID, vmSpec)

	// Wait for VM creation
	mockTask = createMockTask("CREATE_VM", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	vmCreateTask, _ = client.Tasks.Wait(vmCreateTask.ID)

	// Attach disk
	server.SetResponseJson(200, createMockTask("ATTACH_DISK", "QUEUED"))
	attachTask, err := client.VMs.AttachDisk(vmCreateTask.Entity.ID, &VmDiskOperation{DiskID: diskTask.Entity.ID})
	if err != nil {
		t.Error("Not expecting error")
	}
	if attachTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if attachTask.Operation != "ATTACH_DISK" {
		t.Error("Expected task operation to be ATTACH_DISK")
	}
	if attachTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for disk attach
	server.SetResponseJson(200, createMockTask("ATTACH_DISK", "COMPLETED"))
	attachTask, err = client.Tasks.Wait(attachTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if attachTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if attachTask.Operation != "ATTACH_DISK" {
		t.Error("Expected task operation to be ATTACH_DISK")
	}
	if attachTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Detach disk
	server.SetResponseJson(200, createMockTask("DETACH_DISK", "QUEUED"))
	detachTask, err := client.VMs.DetachDisk(vmCreateTask.Entity.ID, &VmDiskOperation{DiskID: diskTask.Entity.ID})
	if err != nil {
		t.Error("Not expecting error")
	}
	if detachTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if detachTask.Operation != "DETACH_DISK" {
		t.Error("Expected task operation to be DETACH_DISK")
	}
	if detachTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for disk detach
	server.SetResponseJson(200, createMockTask("DETACH_DISK", "COMPLETED"))
	detachTask, err = client.Tasks.Wait(detachTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if detachTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if detachTask.Operation != "DETACH_DISK" {
		t.Error("Expected task operation to be DETACH_DISK")
	}
	if detachTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Cleanup disk
	deleteDiskTask, _ := client.Disks.Delete(diskTask.Entity.ID, true)
	server.SetResponseJson(200, createMockTask("DELETE_DISK", "COMPLETED"))
	deleteDiskTask, _ = client.Tasks.Wait(deleteDiskTask.ID)

	// Cleanup VM
	mockTask = createMockTask("DELETE_VM", "QUEUED")
	server.SetResponseJson(200, mockTask)
	vmDeleteTask, err := client.VMs.Delete(vmCreateTask.Entity.ID, true)
	if err != nil {
		t.Error("Not expecting error when deleting VM")
		t.Log(err)
	}

	mockTask = createMockTask("DELETE_VM", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	vmDeleteTask, err = client.Tasks.Wait(vmDeleteTask.ID)
	if err != nil {
		t.Error("Not expecting error when deleting VM")
		t.Log(err)
	}

	// Cleanup image
	imageTask, err = client.Images.Delete(imageTask.Entity.ID)
	imageTask, err = client.Tasks.Wait(imageTask.ID)
	if err != nil {
		t.Error("Not expecting error when deleting image")
		t.Log(err)
	}

	// Cleanup project
	_, err = client.Projects.Delete(projTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting project")
		t.Log(err)
	}

	// Cleanup tenant
	_, err = client.Tenants.Delete(tenantTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting tenant")
		t.Log(err)
	}
}

func TestAttachDetachISO(t *testing.T) {
	if isIntegrationTest() && !isRealAgent() {
		t.Skip("Skipping attach/detach ISO test unless REAL_AGENT env var is set")
	}

	// Create tenant
	mockTask := createMockTask("CREATE_TENANT", "COMPLETED")
	server, client := testSetup()
	server.SetResponseJson(200, mockTask)
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10)}
	tenantTask, _ := client.Tenants.Create(tenantSpec)

	// Create resource ticket
	resSpec := &ResourceTicketCreateSpec{
		Name:   randomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 4, Key: "vm.memory"}},
	}
	_, _ = client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, resSpec)

	// Create project
	projSpec := &ProjectCreateSpec{
		ResourceTicketReservation{
			resSpec.Name,
			[]QuotaLineItem{
				QuotaLineItem{"GB", 4, "vm.memory"},
				QuotaLineItem{"COUNT", 100, "vm"},
			},
		},
		randomString(10),
	}
	projTask, _ := client.Tenants.CreateProject(tenantTask.Entity.ID, projSpec)

	// Create ephemeral disk flavor
	ephemeralFlavor := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "ephemeral-disk.cost"}},
		"ephemeral-disk",
		randomString(10),
	}
	_, _ = client.Flavors.Create(ephemeralFlavor)

	// Create persistent disk flavor
	persistentFlavor := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "persistent-disk.cost"}},
		"persistent-disk",
		randomString(10),
	}
	_, _ = client.Flavors.Create(persistentFlavor)

	// Upload image
	imagePath := "../testdata/tty_tiny.ova"
	imageTask, _ := client.Images.CreateFromFile(imagePath, nil)
	mockTask = createMockTask("CREATE_IMAGE", "COMPLETED")
	server.SetResponseJson(200, mockTask)
	imageTask, _ = client.Tasks.Wait(imageTask.ID)

	vmFlavorSpec := &FlavorCreateSpec{
		Name: randomString(10),
		Kind: "vm",
		Cost: []QuotaLineItem{
			QuotaLineItem{"GB", 2, "vm.memory"},
			QuotaLineItem{"COUNT", 4, "vm.cpu"},
		},
	}
	_, _ = client.Flavors.Create(vmFlavorSpec)

	// Create VM
	mockTask = createMockTask("CREATE_VM", "QUEUED")
	server.SetResponseJson(200, mockTask)
	vmSpec := &VmCreateSpec{
		Flavor:        vmFlavorSpec.Name,
		SourceImageID: imageTask.Entity.ID,
		AttachedDisks: []AttachedDisk{
			AttachedDisk{
				CapacityGB: 1,
				Flavor:     ephemeralFlavor.Name,
				Kind:       "ephemeral-disk",
				Name:       randomString(10),
				State:      "STARTED",
				BootDisk:   true,
			},
		},
		Name: randomString(10),
	}
	vmCreateTask, err := client.Projects.CreateVM(projTask.Entity.ID, vmSpec)

	// Wait for VM creation
	server.SetResponseJson(200, createMockTask("CREATE_VM", "COMPLETED"))
	vmCreateTask, _ = client.Tasks.Wait(vmCreateTask.ID)

	// Attach ISO
	server.SetResponseJson(200, createMockTask("ATTACH_ISO", "COMPLETED"))
	attachIsoTask, err := client.VMs.AttachISO(vmCreateTask.Entity.ID, "../testdata/ttylinux-pc_i486-16.1.iso")
	attachIsoTask, err = client.Tasks.Wait(attachIsoTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if attachIsoTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if attachIsoTask.Operation != "ATTACH_ISO" {
		t.Error("Expected task operation to be ATTACH_ISO")
	}
	if attachIsoTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Detach ISO
	server.SetResponseJson(200, createMockTask("DETACH_ISO", "QUEUED"))
	detachIsoTask, err := client.VMs.DetachISO(vmCreateTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if detachIsoTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if detachIsoTask.Operation != "DETACH_ISO" {
		t.Error("Expected task operation to be DETACH_ISO")
	}
	if detachIsoTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for ISO detach
	server.SetResponseJson(200, createMockTask("DETACH_ISO", "COMPLETED"))
	detachIsoTask, err = client.Tasks.Wait(detachIsoTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if detachIsoTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if detachIsoTask.Operation != "DETACH_ISO" {
		t.Error("Expected task operation to be DETACH_ISO")
	}
	if detachIsoTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Cleanup VM
	server.SetResponseJson(200, createMockTask("DELETE_VM", "QUEUED"))
	vmDeleteTask, err := client.VMs.Delete(vmCreateTask.Entity.ID, true)
	if err != nil {
		t.Error("Not expecting error when deleting VM")
		t.Log(err)
	}

	server.SetResponseJson(200, createMockTask("DELETE_VM", "COMPLETED"))
	vmDeleteTask, err = client.Tasks.Wait(vmDeleteTask.ID)
	if err != nil {
		t.Error("Not expecting error when deleting VM")
		t.Log(err)
	}

	// Cleanup image
	imageTask, err = client.Images.Delete(imageTask.Entity.ID)
	imageTask, err = client.Tasks.Wait(imageTask.ID)
	if err != nil {
		t.Error("Not expecting error when deleting image")
		t.Log(err)
	}

	// Cleanup project
	_, err = client.Projects.Delete(projTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting project")
		t.Log(err)
	}

	// Cleanup tenant
	_, err = client.Tenants.Delete(tenantTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting tenant")
		t.Log(err)
	}
}

func TestVmPowerOnAndOff(t *testing.T) {
	// Create tenant
	server, client := testSetup()
	server.SetResponseJson(200, createMockTask("CREATE_TENANT", "COMPLETED"))
	defer server.Close()

	tenantSpec := &TenantCreateSpec{Name: randomString(10)}
	tenantTask, _ := client.Tenants.Create(tenantSpec)

	// Create resource ticket
	resSpec := &ResourceTicketCreateSpec{
		Name:   randomString(10),
		Limits: []QuotaLineItem{QuotaLineItem{Unit: "GB", Value: 4, Key: "vm.memory"}},
	}
	_, _ = client.Tenants.CreateResourceTicket(tenantTask.Entity.ID, resSpec)

	// Create project
	projSpec := &ProjectCreateSpec{
		ResourceTicketReservation{
			resSpec.Name,
			[]QuotaLineItem{
				QuotaLineItem{"GB", 4, "vm.memory"},
				QuotaLineItem{"COUNT", 100, "vm"},
			},
		},
		randomString(10),
	}
	projTask, _ := client.Tenants.CreateProject(tenantTask.Entity.ID, projSpec)

	// Create flavor
	flavorSpec := &FlavorCreateSpec{
		[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "ephemeral-disk.cost"}},
		"ephemeral-disk",
		randomString(10),
	}
	_, _ = client.Flavors.Create(flavorSpec)

	// Upload image
	imagePath := "../testdata/tty_tiny.ova"
	imageTask, _ := client.Images.CreateFromFile(imagePath, nil)
	server.SetResponseJson(200, createMockTask("CREATE_IMAGE", "COMPLETED"))
	imageTask, _ = client.Tasks.Wait(imageTask.ID)

	vmFlavorSpec := &FlavorCreateSpec{
		Name: randomString(10),
		Kind: "vm",
		Cost: []QuotaLineItem{
			QuotaLineItem{"GB", 2, "vm.memory"},
			QuotaLineItem{"COUNT", 4, "vm.cpu"},
		},
	}
	_, _ = client.Flavors.Create(vmFlavorSpec)

	// Create VM
	server.SetResponseJson(200, createMockTask("CREATE_VM", "QUEUED"))
	vmSpec := &VmCreateSpec{
		Flavor:        vmFlavorSpec.Name,
		SourceImageID: imageTask.Entity.ID,
		AttachedDisks: []AttachedDisk{
			AttachedDisk{
				CapacityGB: 1,
				Flavor:     flavorSpec.Name,
				Kind:       "ephemeral-disk",
				Name:       randomString(10),
				State:      "STARTED",
				BootDisk:   true,
			},
		},
		Name: randomString(10),
	}
	vmCreateTask, err := client.Projects.CreateVM(projTask.Entity.ID, vmSpec)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vmCreateTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if vmCreateTask.Operation != "CREATE_VM" {
		t.Error("Expected task operation to be CREATE_VM")
	}
	if vmCreateTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for VM creation
	server.SetResponseJson(200, createMockTask("CREATE_VM", "COMPLETED"))
	vmCreateTask, err = client.Tasks.Wait(vmCreateTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vmCreateTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Power on
	server.SetResponseJson(200, createMockTask("START_VM", "QUEUED"))
	powerOnTask, err := client.VMs.Operation(vmCreateTask.Entity.ID, &VmOperation{Operation: "START_VM"})
	if err != nil {
		t.Error("Not expecting error")
	}
	if powerOnTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if powerOnTask.Operation != "START_VM" {
		t.Error("Expected task operation to be START_VM")
	}
	if powerOnTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for power on
	server.SetResponseJson(200, createMockTask("START_VM", "COMPLETED"))
	powerOnTask, err = client.Tasks.Wait(powerOnTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if powerOnTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if powerOnTask.Operation != "START_VM" {
		t.Error("Expected task operation to be START_VM")
	}
	if powerOnTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Power off
	server.SetResponseJson(200, createMockTask("STOP_VM", "QUEUED"))
	powerOffTask, err := client.VMs.Operation(vmCreateTask.Entity.ID, &VmOperation{Operation: "STOP_VM"})
	if err != nil {
		t.Error("Not expecting error")
	}
	if powerOffTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if powerOffTask.Operation != "STOP_VM" {
		t.Error("Expected task operation to be STOP_VM")
	}
	if powerOffTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	// Wait for power off
	server.SetResponseJson(200, createMockTask("STOP_VM", "COMPLETED"))
	powerOffTask, err = client.Tasks.Wait(powerOffTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if powerOffTask == nil {
		t.Error("Not expecting task to be nil")
	}
	if powerOffTask.Operation != "STOP_VM" {
		t.Error("Expected task operation to be STOP_VM")
	}
	if powerOffTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Cleanup VM
	server.SetResponseJson(200, createMockTask("DELETE_VM", "QUEUED"))
	vmDeleteTask, err := client.VMs.Delete(vmCreateTask.Entity.ID, true)
	if err != nil {
		t.Error("Not expecting error")
		t.Log(err)
	}
	if vmDeleteTask.Operation != "DELETE_VM" {
		t.Error("Expected task operation to be DELETE_VM")
	}
	if vmDeleteTask.State != "QUEUED" {
		t.Error("Expected task status to be QUEUED")
	}

	server.SetResponseJson(200, createMockTask("DELETE_VM", "COMPLETED"))
	vmDeleteTask, err = client.Tasks.Wait(vmDeleteTask.ID)
	if err != nil {
		t.Error("Not expecting error")
	}
	if vmDeleteTask.Operation != "DELETE_VM" {
		t.Error("Expected task operation to be DELETE_VM")
	}
	if vmDeleteTask.State != "COMPLETED" {
		t.Error("Expected task status to be COMPLETED")
	}

	// Cleanup image
	imageTask, err = client.Images.Delete(imageTask.Entity.ID)
	imageTask, err = client.Tasks.Wait(imageTask.ID)
	if err != nil {
		t.Error("Not expecting error when deleting image")
		t.Log(err)
	}

	// Cleanup project
	_, err = client.Projects.Delete(projTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting project")
		t.Log(err)
	}

	// Cleanup tenant
	_, err = client.Tenants.Delete(tenantTask.Entity.ID)
	if err != nil {
		t.Error("Not expecting error when deleting tenant")
		t.Log(err)
	}
}
