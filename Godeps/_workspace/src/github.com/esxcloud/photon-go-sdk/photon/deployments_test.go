package photon

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Deployment", func() {
	var (
		server         *testServer
		client         *Client
		deploymentSpec *DeploymentCreateSpec
	)

	BeforeEach(func() {
		if isIntegrationTest() {
			Skip("Skipping deployment test on integration mode. Need undeployed environment")
		}
		server, client = testSetup()
		deploymentSpec = &DeploymentCreateSpec{
			ImageDatastores: []string{randomString(10, "go-sdk-deployment-")},
			Auth: &AuthInfo{
				Enabled: false,
			},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("CreateGetAndDeleteDeployment", func() {
		It("Deployment create and delete succeeds", func() {
			if isIntegrationTest() {
				Skip("Skipping deployment test on integration mode. Need undeployed environment")
			}
			mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.Create(deploymentSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockDeployment := Deployment{
				ImageDatastores: deploymentSpec.ImageDatastores,
				Auth:            &AuthInfo{Enabled: false},
			}
			server.SetResponseJson(200, mockDeployment)
			deployment, err := client.Deployments.Get(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(deployment).ShouldNot(BeNil())
			Expect(deployment.ImageDatastores).Should(Equal(deploymentSpec.ImageDatastores))
			Expect(deployment.ID).Should(Equal(task.Entity.ID))

			server.SetResponseJson(200, &Deployments{[]Deployment{mockDeployment}})
			deploymentList, err := client.Deployments.GetAll()
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(deploymentList).ShouldNot(BeNil())

			var found bool
			for _, d := range deploymentList.Items {
				if reflect.DeepEqual(d.ImageDatastores, deploymentSpec.ImageDatastores) {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Deployments.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("GetDeployment", func() {
		It("Get Hosts succeeds", func() {
			mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.Create(deploymentSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			mockTask = createMockTask("CREATE_HOST", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			hostSpec := &HostCreateSpec{
				Username: randomString(10),
				Password: randomString(10),
				Address:  randomAddress(),
				Tags:     []string{"CLOUD"},
				Metadata: map[string]string{"test": "go-sdk-host"},
			}

			mockTask = createMockTask("CREATE_HOST", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			hostTask, err := client.Hosts.Create(hostSpec, task.Entity.ID)
			hostTask, err = client.Tasks.Wait(hostTask.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &Hosts{[]Host{Host{Tags: hostSpec.Tags, ID: task.Entity.ID}}})
			hostList, err := client.Deployments.GetHosts(task.Entity.ID)

			var found bool
			for _, host := range hostList.Items {
				if host.ID == hostTask.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Deployments.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})

		Context("GetVM needs a vm", func() {
			var (
				tenantID     string
				resName      string
				projID       string
				imageID      string
				flavorSpec   *FlavorCreateSpec
				vmFlavorSpec *FlavorCreateSpec
				vmSpec       *VmCreateSpec
			)

			BeforeEach(func() {
				tenantID = createTenant(server, client)
				resName = createResTicket(server, client, tenantID)
				projID = createProject(server, client, tenantID, resName)
				imageID = createImage(server, client)
				flavorSpec = &FlavorCreateSpec{
					[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "ephemeral-disk.cost"}},
					"ephemeral-disk",
					randomString(10, "go-sdk-flavor-"),
				}

				_, err := client.Flavors.Create(flavorSpec)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				vmFlavorSpec = &FlavorCreateSpec{
					Name: randomString(10, "go-sdk-flavor-"),
					Kind: "vm",
					Cost: []QuotaLineItem{
						QuotaLineItem{"GB", 2, "vm.memory"},
						QuotaLineItem{"COUNT", 4, "vm.cpu"},
					},
				}
				_, err = client.Flavors.Create(vmFlavorSpec)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				vmSpec = &VmCreateSpec{
					Flavor:        vmFlavorSpec.Name,
					SourceImageID: imageID,
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
					Name: randomString(10, "go-sdk-vm-"),
				}
			})

			AfterEach(func() {
				cleanVMs(client, projID)
				cleanImages(client)
				cleanFlavors(client)
				cleanTenants(client)
			})

			It("Get VMs succeeds", func() {
				mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
				server.SetResponseJson(200, mockTask)

				task, err := client.Deployments.Create(deploymentSpec)
				task, err = client.Tasks.Wait(task.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				mockTask = createMockTask("CREATE_VM", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				vmTask, err := client.Projects.CreateVM(projID, vmSpec)
				vmTask, err = client.Tasks.Wait(vmTask.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				server.SetResponseJson(200, &VMs{[]VM{VM{Name: vmSpec.Name}}})
				vmList, err := client.Deployments.GetVms(task.Entity.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				var found bool
				for _, vm := range vmList.Items {
					if vm.Name == vmSpec.Name && vm.ID == vmTask.Entity.ID {
						found = true
						break
					}
				}
				Expect(found).Should(BeTrue())

				mockTask = createMockTask("DELETE_VM", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				task, err = client.VMs.Delete(task.Entity.ID)
				task, err = client.Tasks.Wait(task.ID)

				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				mockTask = createMockTask("DELETE_DEPLOYMENT", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				task, err = client.Deployments.Delete(task.Entity.ID)
				task, err = client.Tasks.Wait(task.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
			})
		})

	})

	Describe("DeployAndDestroyDeployment", func() {
		It("Deploy and Destroy a deployment succeeds", func() {
			if isIntegrationTest() {
				Skip("Skipping deployment test on integration mode. Need undeployed environment")
			}
			mockTask := createMockTask("PERFORM_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.Deploy("deployment-ID")
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("PERFORM_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("PERFORM_DELETE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Deployments.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("PERFORM_DELETE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("InitializeAndFinalizeMigrateDeployment", func() {
		It("Initialize and Finalize deployment migration succeeds", func() {
			if isIntegrationTest() {
				Skip("Skipping deployment test on integration mode. Need undeployed environment")
			}
			mockTask := createMockTask("PERFORM_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.Deploy("deployment-ID")
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("PERFORM_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("INITIALIZE_MIGRATE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err = client.Deployments.InitializeDeploymentMigration("sourceAddr", "deployment-ID")
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("INITIALIZE_MIGRATE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("FINALIZE_MIGRATE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err = client.Deployments.FinalizeDeploymentMigration("sourceAddr", "deployment-ID")
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("FINALIZE_MIGRATE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("PERFORM_DELETE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Deployments.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("PERFORM_DELETE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})
})
