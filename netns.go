package netns

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var (
	NamespaceCloseError = errors.New("namespace close")
)

type NetworkNsHandler struct {
	ns     *os.File
	closed bool
}

type NetworkStats []NetworkStat

type networkInfo struct {
	Bytes      int64
	Packets    int64
	Drop       int64
	Errs       int64
	Fifo       int64
	Frame      int64
	Compressed int64
	Multicast  int64
}

type NetworkStat struct {
	Interface string
	Received  networkInfo
	Transmit  networkInfo
}

func Setns(pid string) (*NetworkNsHandler, error) {
	var err error
	h := &NetworkNsHandler{}
	h.ns, err = os.Open("/proc/" + pid + "/ns/net")
	if err != nil {
		return nil, err
	}

	err = setns(h.ns.Fd())
	if err != nil {
		h.ns.Close()
		return nil, err
	}

	return h, nil
}

func setns(fd uintptr) error {
	_, _, err := syscall.Syscall(308, fd, uintptr(0), uintptr(0))
	if uintptr(err) == 0 {
		return nil
	}
	return err
}

func (h *NetworkNsHandler) Close() error {
	h.closed = true
	return h.ns.Close()
}

func (h *NetworkNsHandler) Stats() (NetworkStats, error) {
	if h.closed {
		return nil, NamespaceCloseError
	}
	netStatsFile, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer netStatsFile.Close()

	var stats NetworkStats
	reader := bufio.NewReader(netStatsFile)

	// Pass the header
	// Inter-|   Receive                                                |  Transmit
	//  face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
	reader.ReadString('\n')
	reader.ReadString('\n')

	var line string
	for err == nil {
		line, err = reader.ReadString('\n')
		if line == "" {
			continue
		}
		stats = append(stats, buildNetworkStat(line))
	}
	return stats, nil
}

func buildNetworkStat(line string) NetworkStat {
	fields := strings.Fields(line)
	interfaceName := strings.TrimSuffix(fields[0], ":")
	return NetworkStat{
		Interface: interfaceName,
		Received:  toNetworkInfo(fields[1:9]),
		Transmit:  toNetworkInfo(fields[9:17]),
	}
}

func toNetworkInfo(fields []string) networkInfo {
	return networkInfo{
		Bytes:      toInt(fields[0]),
		Packets:    toInt(fields[1]),
		Errs:       toInt(fields[2]),
		Drop:       toInt(fields[3]),
		Fifo:       toInt(fields[4]),
		Frame:      toInt(fields[5]),
		Compressed: toInt(fields[6]),
		Multicast:  toInt(fields[7]),
	}
}

func toInt(str string) int64 {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}
