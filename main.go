package main

import(
	"fmt"
	"log"
	"os"
	"bufio"
	"time"
	"sync"
	"strings"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket"
)

func main() {
	var once sync.Once
	once.Do(func(){dviceName()})
	fmt.Println("キャプチャーするデバイスとポートを入力してください")

	//標準入力の受け取り
	stadin := bufio.NewScanner(os.Stdin)
	stadin.Scan()
	name := strings.Fields(stadin.Text())

	//キャプチャ開始
	capture(name[0], name[1])
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

//パケットをキャプチャしてくれる関数
func capture(name string, port string){
	var (

		//ネットワークデバイス名は環境による
    divice_name  string = name
    snapshot_len int32  = 1024
    promiscuous  bool   = false
    err          error
    timeout      time.Duration = 30 * time.Second
    handle       *pcap.Handle
)

	handle, err = pcap.OpenLive(divice_name, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	//tcpポートのフィルター
	var filter string = "tcp and port " + port
	fmt.Print(filter)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Only capturing TCP port " + port + " packets.")

	//パケットソースの処理
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}
}