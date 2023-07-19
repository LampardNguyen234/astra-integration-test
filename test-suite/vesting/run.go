package vesting

func (s *VestingSuite) RunTest() {
	s.Start()
	//s.runCreateVestingTests()
	s.runClawBackVestingTest()
	s.Finished()
}
