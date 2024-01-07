package resource

import (
        "strconv"
        "strings"
)

func Split_ip_with_route_domain(address string) (ip string, rd string) {
        // Split the address into the ip and routeDomain (optional) parts
        //     address is of the form: <ipv4_or_ipv6>[%<routeDomainID>]
        match := strings.Split(address, "%")
        if len(match) == 2 {
                _, err := strconv.Atoi(match[1])
                //Matches only when RD contains number, Not allowing RD has 80f
                if err == nil {
                        ip = match[0]
                        rd = match[1]
                } else {
                        ip = address
                }
        } else {
                ip = match[0]
        }
        return
}
