package esxcloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Disk", func() {
	var (
		server     *testServer
		client     *Client
		tenantID   string
		resName    string
		projID     string
		flavorName string
		flavorID   string
		diskSpec   *DiskCreateSpec
	)

	BeforeEach(func() {
		server, client = testSetup()
		tenantID = createTenant(server, client)
		resName = createResTicket(server, client, tenantID)
		projID = createProject(server, client, tenantID, resName)
		flavorName, flavorID = createFlavor(server, client)
		diskSpec = &DiskCreateSpec{
			Flavor:     flavorName,
			Kind:       "persistent-disk",
			CapacityGB: 2,
			Name:       randomString(10, "go-sdk-disk-"),
		}
	})

	AfterEach(func() {
		cleanDisks(client, projID)
		cleanFlavors(client)
		cleanTenants(client)
		server.Close()
	})

	Describe("CreateAndDeleteDisk", func() {
		It("Disk create and delete succeeds", func() {
			mockTask := createMockTask("CREATE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Projects.CreateDisk(projID, diskSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_DISK"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("DELETE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Disks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_DISK"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("GetDisk", func() {
		It("Get disk returns a disk ID", func() {
			mockTask := createMockTask("CREATE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Projects.CreateDisk(projID, diskSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			diskMock := &PersistentDisk{
				Name:       diskSpec.Name,
				Flavor:     diskSpec.Flavor,
				CapacityGB: diskSpec.CapacityGB,
				Kind:       diskSpec.Kind,
			}
			server.SetResponseJson(200, diskMock)
			disk, err := client.Disks.Get(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(disk.Name).Should(Equal(diskSpec.Name))
			Expect(disk.Flavor).Should(Equal(diskSpec.Flavor))
			Expect(disk.Kind).Should(Equal(diskSpec.Kind))
			Expect(disk.CapacityGB).Should(Equal(diskSpec.CapacityGB))
			Expect(disk.ID).Should(Equal(task.Entity.ID))

			mockTask = createMockTask("DELETE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Disks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})

	Describe("GetTasks", func() {
		It("GetTasks returns a completed task", func() {
			mockTask := createMockTask("CREATE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Projects.CreateDisk(projID, diskSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
			taskList, err := client.Disks.GetTasks(task.Entity.ID, &TaskGetOptions{})

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(taskList).ShouldNot(BeNil())
			Expect(taskList.Items).Should(ContainElement(*task))

			mockTask = createMockTask("DELETE_DISK", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Disks.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})
})
