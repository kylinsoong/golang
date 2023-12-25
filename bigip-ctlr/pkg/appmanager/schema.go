package appmanager

import (
    "net"
    "strconv"
    "github.com/xeipuuv/gojsonschema"
    "github.com/kylinsoong/bigip-ctlr/pkg/resource"
)

type BigIPv4FormatChecker struct{}

func (f BigIPv4FormatChecker) IsFormat(input interface{}) bool {
        var strInput = input.(string)
        ip, rd := resource.Split_ip_with_route_domain(strInput)
        if rd != "" {
                if _, err := strconv.Atoi(rd); err != nil {
                        return false
                }
        }

        address := net.ParseIP(ip)
        if nil == address.To4() {
                return false
        }
        return true
}

type BigIPv6FormatChecker struct{}

func (f BigIPv6FormatChecker) IsFormat(input interface{}) bool {
        var strInput = input.(string)
        ip, rd := resource.Split_ip_with_route_domain(strInput)
        if rd != "" {
                if _, err := strconv.Atoi(rd); err != nil {
                        return false
                }
        }

        address := net.ParseIP(ip)
        if nil == address.To16() {
                return false
        }

        return true
}

func RegisterBigIPSchemaTypes() {
        gojsonschema.FormatCheckers.Add("bigipv4", BigIPv4FormatChecker{})
        gojsonschema.FormatCheckers.Add("bigipv6", BigIPv6FormatChecker{})
}
