package redis

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClusterClient", func() {
	var subject *ClusterClient

	var populate = func() {
		subject.setSlots([]ClusterSlotInfo{
			{0, 4095, []string{"127.0.0.1:7000", "127.0.0.1:7004"}},
			{12288, 16383, []string{"127.0.0.1:7003", "127.0.0.1:7007"}},
			{4096, 8191, []string{"127.0.0.1:7001", "127.0.0.1:7005"}},
			{8192, 12287, []string{"127.0.0.1:7002", "127.0.0.1:7006"}},
		})
	}

	BeforeEach(func() {
		var err error
		subject = NewClusterClient(&ClusterOptions{
			Addrs: []string{"127.0.0.1:6379", "127.0.0.1:7003", "127.0.0.1:7006"},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		subject.Close()
	})

	It("should initialize", func() {
		Expect(subject.addrs).To(HaveLen(3))
		Expect(subject._reload).To(Equal(uint32(1)))
	})

	It("should update slots cache", func() {
		populate()
		Expect(subject.slots[0]).To(Equal([]string{"127.0.0.1:7000", "127.0.0.1:7004"}))
		Expect(subject.slots[4095]).To(Equal([]string{"127.0.0.1:7000", "127.0.0.1:7004"}))
		Expect(subject.slots[4096]).To(Equal([]string{"127.0.0.1:7001", "127.0.0.1:7005"}))
		Expect(subject.slots[8191]).To(Equal([]string{"127.0.0.1:7001", "127.0.0.1:7005"}))
		Expect(subject.slots[8192]).To(Equal([]string{"127.0.0.1:7002", "127.0.0.1:7006"}))
		Expect(subject.slots[12287]).To(Equal([]string{"127.0.0.1:7002", "127.0.0.1:7006"}))
		Expect(subject.slots[12288]).To(Equal([]string{"127.0.0.1:7003", "127.0.0.1:7007"}))
		Expect(subject.slots[16383]).To(Equal([]string{"127.0.0.1:7003", "127.0.0.1:7007"}))
		Expect(subject.addrs).To(Equal([]string{
			"127.0.0.1:6379",
			"127.0.0.1:7003",
			"127.0.0.1:7006",
			"127.0.0.1:7000",
			"127.0.0.1:7004",
			"127.0.0.1:7007",
			"127.0.0.1:7001",
			"127.0.0.1:7005",
			"127.0.0.1:7002",
		}))
	})

	It("should check if reload is due", func() {
		subject._reload = 0
		Expect(subject._reload).To(Equal(uint32(0)))
		subject.scheduleReload()
		Expect(subject._reload).To(Equal(uint32(1)))
	})
})