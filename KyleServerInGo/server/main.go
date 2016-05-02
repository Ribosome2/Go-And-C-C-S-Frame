//服务端解包过程
package main
import (
	"../protocol"
	"fmt"
	"net"
	"os"
)
type User struct {
	Name string
	ID int
	Initiated bool
	UChannel chan []byte
	Connection *net.Conn
}
//负责从协议通道里面处理收到的数据， 因为解包的时候已经做了粘包的处理
//这里每次接受到的应该都是一个整包
func (u User) listenForNetPacket(readerChannel chan []byte)() {
	for {
		select {
		case data := <-readerChannel:
			ParsePacket(u,data)
			Log(string(data))
		}
	}
}
func (u User) handleConnection(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go u.listenForNetPacket(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		//解码之后产生新的缓存区内容
		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}
var Users []*User

func main() {
	netListen, err := net.Listen("tcp", ":9988")
	CheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		user := User{Name:"anonymous",ID:0,Initiated:false}
		Users = append(Users, &user)
		Log(conn.RemoteAddr().String(), " tcp connect success")
		go user.handleConnection(conn)
	}
}
func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
