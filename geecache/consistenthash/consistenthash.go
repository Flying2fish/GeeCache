package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(b []byte) uint32

type Map struct {
	hash     Hash
	replicas int //一个物理节点的虚拟节点的个数
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil { //如果没有指定，则使用crc32.ChecksumIEEE初始化
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) { //创建物理节点
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hashvalue := int(m.hash([]byte(strconv.Itoa(i) + key))) //获取key这个物理节点第i个虚拟节点的哈希值
			m.keys = append(m.keys, hashvalue)
			m.hashMap[hashvalue] = key //虚拟节点的哈希值输入得到物理节点名称
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string { //输入key得到从哪个真实物理节点获得数据
	if len(m.keys) == 0 {
		return ""
	}

	hashvalue := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hashvalue
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
