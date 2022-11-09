package enums

type Node struct {
	Host string
	Test bool
}

func CreateNode(host string, isTest bool) Node {
	return Node{
		Host: host,
		Test: isTest,
	}
}

var (
	MAIN_NODE = Node{
		Host: "https://btc1.trezor.io",
		Test: false,
	}
	CUSTOM_NODE = Node{
		Host: "http://49.12.103.5:9130",
		Test: false,
	}
)
