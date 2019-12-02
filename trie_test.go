package trie_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/dghubble/trie"
	"github.com/yl2chen/cidranger"
)

var pathKeys []string
var cidrs []string

var tree = trie.NewPathTrie()
var ranger = cidranger.NewPCTrieRanger()

func init() {
	// path keys
	data, _ := IPRangeFromFile("xxx")
	pathKeys = make([]string, len(data))
	cidrs = make([]string, len(data))

	i := 0
	for cidr, name := range data {
		ip, network, _ := net.ParseCIDR(cidr)
		key := fmt.Sprintf("%b", ip)
		pathKeys[i] = key
		cidrs[i] = cidr

		tree.Put(key, name)
		_ = ranger.Insert(NewRangerEntry(*network, name))
		i++
	}

}

func BenchmarkPathTriePutPathKey(b *testing.B) {
	putree := trie.NewPathTrie()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		putree.Put(pathKeys[i%len(pathKeys)], i)
	}
}

func BenchmarkPathTrieGetPathKey(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tree.Get(pathKeys[i%len(pathKeys)])
	}
}

func BenchmarkCIDRangerGetKey(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ranger.ContainingNetworks(net.ParseIP(pathKeys[i%len(pathKeys)]))
	}
}

func BenchmarkCIDRangerPutKey(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, network, _ := net.ParseCIDR(cidrs[i%len(cidrs)])
		_ = ranger.Insert(NewRangerEntry(*network, ""))
	}
}
