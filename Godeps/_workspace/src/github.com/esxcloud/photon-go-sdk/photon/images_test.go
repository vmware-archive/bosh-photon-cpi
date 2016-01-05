package photon

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image", func() {
	var (
		server *testServer
		client *Client
	)

	BeforeEach(func() {
		server, client = testSetup()
	})

	AfterEach(func() {
		cleanImages(client)
		server.Close()
	})

	Describe("CreateAndDeleteImage", func() {
		It("Image create from file and delete succeeds", func() {
			mockTask := createMockTask("CREATE_IMAGE", "COMPLETED", createMockStep("UPLOAD_IMAGE", "COMPLETED"))
			server.SetResponseJson(200, mockTask)

			// create image from file
			imagePath := "../testdata/tty_tiny.ova"
			task, err := client.Images.CreateFromFile(imagePath, &ImageCreateOptions{ReplicationType: "ON_DEMAND"})
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_IMAGE"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("DELETE_IMAGE", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Images.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_IMAGE"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})

		It("Image create and delete succeeds", func() {
			mockTask := createMockTask("CREATE_IMAGE", "COMPLETED", createMockStep("UPLOAD_IMAGE", "COMPLETED"))
			server.SetResponseJson(200, mockTask)

			imagePath := "../testdata/tty_tiny.ova"
			file, err := os.Open(imagePath)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			task, err := client.Images.Create(file, "tty_tiny.ova", &ImageCreateOptions{ReplicationType: "ON_DEMAND"})
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_IMAGE"))
			Expect(task.State).Should(Equal("COMPLETED"))

			err = file.Close()
			Expect(err).Should(BeNil())

			mockTask = createMockTask("DELETE_IMAGE", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Images.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_IMAGE"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("GetImage", func() {
		It("Get image succeeds", func() {
			mockTask := createMockTask("CREATE_IMAGE", "COMPLETED", createMockStep("UPLOAD_IMAGE", "COMPLETED"))
			server.SetResponseJson(200, mockTask)

			// create image from file
			imagePath := "../testdata/tty_tiny.ova"
			task, err := client.Images.CreateFromFile(imagePath, &ImageCreateOptions{ReplicationType: "ON_DEMAND"})
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, Image{Name: "tty_tiny.ova"})
			image, err := client.Images.Get(task.Entity.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(image).ShouldNot(BeNil())
			Expect(image.ID).Should(Equal(task.Entity.ID))
			Expect(image.Name).Should(Equal("tty_tiny.ova"))

			mockTask = createMockTask("DELETE_IMAGE", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Images.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})

		It("GetAll image succeeds", func() {
			mockTask := createMockTask("CREATE_IMAGE", "COMPLETED", createMockStep("UPLOAD_IMAGE", "COMPLETED"))
			server.SetResponseJson(200, mockTask)

			// create image from file
			imagePath := "../testdata/tty_tiny.ova"
			task, err := client.Images.CreateFromFile(imagePath, &ImageCreateOptions{ReplicationType: "ON_DEMAND"})
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &Images{[]Image{Image{Name: "tty_tiny.ova"}}})
			imageList, err := client.Images.GetAll()
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(imageList).ShouldNot(BeNil())

			var found bool
			for _, image := range imageList.Items {
				if image.Name == "tty_tiny.ova" {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())

			mockTask = createMockTask("DELETE_IMAGE", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Images.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})

		It("GetTasks returns a completed task", func() {
			mockTask := createMockTask("CREATE_IMAGE", "COMPLETED", createMockStep("UPLOAD_IMAGE", "COMPLETED"))
			server.SetResponseJson(200, mockTask)

			// create image from file
			imagePath := "../testdata/tty_tiny.ova"
			task, err := client.Images.CreateFromFile(imagePath, &ImageCreateOptions{ReplicationType: "ON_DEMAND"})
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, &TaskList{[]Task{*mockTask}})
			taskList, err := client.Images.GetTasks(task.Entity.ID, &TaskGetOptions{})

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(taskList).ShouldNot(BeNil())
			Expect(taskList.Items).Should(ContainElement(*task))

			mockTask = createMockTask("DELETE_IMAGE", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Images.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
		})
	})
})