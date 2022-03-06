/*
Copyright © 2022
Author Bhakiyaraj Kalimuthu
Email bhakiya.kalimuthu@gmail.com
*/

package internal

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func randomIPFromRange() (net.IP, error) {
GENERATE:
	ip, ipnet, err := net.ParseCIDR("192.168.100.0/24")
	if err != nil {
		return nil, err
	}

	// The number of leading 1s in the mask
	ones, _ := ipnet.Mask.Size()
	quotient := ones / 8
	remainder := ones % 8

	// create random 4-byte byte slice
	r := make([]byte, 4)
	rand.Read(r)

	for i := 0; i <= quotient; i++ {
		if i == quotient {
			shifted := byte(r[i]) >> remainder
			r[i] = ^ipnet.IP[i] & shifted
		} else {
			r[i] = ipnet.IP[i]
		}
	}
	ip = net.IPv4(r[0], r[1], r[2], r[3])

	if ip.Equal(ipnet.IP) /*|| ip.Equal(broadcast) */ {
		// we got unlucky. The host portion of our ipv4 address was
		// either all 0s (the network address) or all 1s (the broadcast address)
		goto GENERATE
	}
	return ip, nil
}

func getIP() string {
	netIP, err := randomIPFromRange()
	if err != nil {
		return fmt.Sprintf("3.22.15.165, 172.70.38.%d", randInt(1, 99))
	} else {
		return netIP.String()
	}
}

func randInt(min int, max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	val := r.Intn(max - min + 1)
	fmt.Println(val)
	return val
}
