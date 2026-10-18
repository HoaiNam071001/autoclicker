package engine

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

// SNTP tối giản (RFC 4330) cho node Chờ đến giờ canh giờ chuẩn (bước 35 N6): độ lệch giờ chuẩn − giờ máy.
var ntpServers = []string{"time.google.com:123", "time.cloudflare.com:123", "pool.ntp.org:123"}

// ntpEpoch – 1900-01-01, gốc thời gian của NTP.
var ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

// ntpOffset hỏi lần lượt các server tới khi một server trả lời (mỗi server tối đa 2 s).
func ntpOffset(ctx context.Context) (time.Duration, error) {
	err := errors.New("không có server giờ chuẩn")
	for _, srv := range ntpServers {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		var off time.Duration
		if off, err = ntpQuery(srv, 2*time.Second); err == nil {
			return off, nil
		}
	}
	return 0, err
}

func ntpQuery(addr string, timeout time.Duration) (time.Duration, error) {
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))
	req := make([]byte, 48)
	req[0] = 0x23 // LI 0, phiên bản 4, chế độ client
	t1 := time.Now()
	if _, err := conn.Write(req); err != nil {
		return 0, err
	}
	resp := make([]byte, 48)
	if _, err := conn.Read(resp); err != nil {
		return 0, err
	}
	t4 := time.Now()
	if resp[0]&0x07 != 4 || resp[1] == 0 { // không phải server / stratum 0 (kiss-o'-death)
		return 0, errors.New("server giờ chuẩn trả lời không hợp lệ")
	}
	t2, t3 := ntpTime(resp[32:40]), ntpTime(resp[40:48])
	return (t2.Sub(t1) + t3.Sub(t4)) / 2, nil
}

// ntpTime – 64 bit: giây (32) + phần giây (32) từ ntpEpoch.
func ntpTime(b []byte) time.Time {
	sec, frac := binary.BigEndian.Uint32(b[:4]), binary.BigEndian.Uint32(b[4:])
	return ntpEpoch.Add(time.Duration(sec)*time.Second + time.Duration(uint64(frac)*uint64(time.Second)>>32))
}
