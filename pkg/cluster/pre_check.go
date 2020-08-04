package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Fize/n/pkg/utils"
)

// PreCheck Inspection of legality
func PreCheck(podcidr, svccidr, svcportrange string, m ...string) error {
	if err := numOfMaster(m...); err != nil {
		return err
	}
	if err := validateCIDR(podcidr); err != nil {
		return fmt.Errorf("pod network cidr: %v", err)
	}
	if err := validateCIDR(svccidr); err != nil {
		return fmt.Errorf("service cidr: %v", err)
	}
	if err := validateRange(svcportrange); err != nil {
		return err
	}
	return nil
}

// Determines whether the number of masters is odd or even,
// if even return a error.
func numOfMaster(masters ...string) error {
	if len(masters)%2 == 0 {
		return fmt.Errorf("the number of master can only be odd")
	}
	return nil
}

func validateCIDR(cidr string) error {
	if cidr == "" {
		return nil
	}
	cidrSlice := strings.Split(cidr, "/")
	if len(cidrSlice) != 2 {
		return fmt.Errorf("invalid cidr")
	}
	segment := cidrSlice[0]
	netmask, err := strconv.Atoi(cidrSlice[1])
	if err != nil {
		return err
	}
	if netmask >= 32 || netmask < 8 {
		return fmt.Errorf("invalid netmask")
	}
	if err := utils.ValidataIP(segment); err != nil {
		return err
	}
	return nil
}

func validateRange(scope string) error {
	if scope == "" {
		return nil
	}
	scopeSlice := strings.Split(scope, "-")
	if len(scopeSlice) != 2 {
		return fmt.Errorf("invalid scope")
	}
	start, err := strconv.Atoi(scopeSlice[0])
	if err != nil {
		return err
	}
	end, err := strconv.Atoi(scopeSlice[1])
	if err != nil {
		return err
	}
	if start <= 0 || end > 65535 {
		return fmt.Errorf("invalid range")
	}
	return nil
}
