package enums

type Node struct {
	Http string
	Ws   string
	Test bool
}

func CreateNode(host string, isTest bool) Node {
	return Node{
		Http: host,
		Test: isTest,
	}
}

var (
	MAIN_NODE = Node{
		Http: "https://btc1.trezor.io",
		Ws:   "wss://btc1.trezor.io",
		Test: false,
	}
	CUSTOM_NODE = Node{
		Http: "http://49.12.103.5:9130",
		Ws:   "ws://49.12.103.5:9130",
		Test: false,
	}
	TEST_NODE = Node{
		Http: "https://tbtc1.trezor.io",
		Ws:   "wss://tbtc1.trezor.io",
		Test: true,
	}
)
