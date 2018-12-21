package main

import(
	"fmt"
	"log"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket"
	"time"
)

var(
	//ネットワークデバイス名は環境による
	divice_name string = "enp0s3"

	snapshot_len int32 = 1024
	promiscuous  bool = false
	err error
	timeout time.Duration = 30 * time.Second
	handle *pcap.Handle
)

func main() {
	handle, err = pcap.OpenLive(divice_name, snapshot_len, promiscuous, timeout)

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	//tcpポートのフィルター
	var filter string = "tcp and port 25565"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Only capturing TCP port 25565 packets.")

	//パケットソースの処理
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
			fmt.Println(packet)
	}

}

//ネットワークデバイスの名前を取得する関数
func dviceName(){
	devices, err := pcap.FindAllDevs()
	if err != nil {
			log.Fatal(err)
	}

	fmt.Println("Devices found:")
	//ネットワークデバイスの名前を取得
	for _, device := range devices {
			fmt.Printf("%s\n", device.Name)
	}

}