package vesting

import (
	. "github.com/LampardNguyen234/astra-integration-test/framework"
)

func (s *VestingSuite) RunTest() {
	s.Start()
	root := Describe(
		s.Name(),
		s.testCreateVesting(),
		s.testClawBackVesting(),
	)
	root.Run()
	root.Report()
	s.Finished()
}
